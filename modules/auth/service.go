package auth

import (
	"errors"
	"log"
	"os"
	"stable/database/entities"
	"stable/modules/users"
	"stable/packages/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input RegisterInput) (entities.User, error)
	Login(input LoginInput) (string, entities.User, error)
	ResendVerificationCode(email string) error
}

type authService struct {
	repo users.Repository
}

func NewService(repo users.Repository) Service {
	return &authService{repo: repo}
}

// ── REGISTER ─────────────────────────────────────────────────
func (s *authService) Register(input RegisterInput) (entities.User, error) {
	user := entities.User{
		Username:     input.Username,
		Email:        input.Email,
		AuthProvider: "Password",
		Role:         "user",
		IsVerified:   false,
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(hashedPassword)

	// Cek email sudah terdaftar
	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return user, errors.New("email already registered")
	}

	// Buat kode verifikasi
	code := utils.GenerateVerificationCode()
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	hashedCodeStr := string(hashedCode)
	user.VerificationCode = &hashedCodeStr

	// Simpan user ke database
	if err := s.repo.Create(&user); err != nil {
		return user, err
	}

	// Kirim email verifikasi (async)
	go func() {
		emailBody := "Your verification code is: " + code
		if err := utils.SendEmail(user.Email, "Email Verification", emailBody); err != nil {
			log.Println("[EMAIL] FAILED:", err)
		}
	}()

	return user, nil
}

// ── LOGIN ─────────────────────────────────────────────────────
func (s *authService) Login(input LoginInput) (string, entities.User, error) {
	// Cari user berdasarkan email
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", entities.User{}, errors.New("email not found")
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", entities.User{}, errors.New("invalid password")
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		return "", entities.User{}, err
	}

	return token, user, nil
}

// ── RESEND VERIFICATION ───────────────────────────────────────
func (s *authService) ResendVerificationCode(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	code := utils.GenerateVerificationCode()
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedCodeStr := string(hashedCode)
	user.VerificationCode = &hashedCodeStr

	if err := s.repo.Update(user); err != nil {
		return err
	}

	emailBody := `
		<h2>Email Verification</h2>
		<p>Your verification code is: <strong>` + code + `</strong></p>
		<p>This code will expire in 15 minutes.</p>
	`
	go utils.SendEmail(user.Email, "Email Verification", emailBody)

	return nil
}

// ── GENERATE JWT TOKEN (internal) ────────────────────────────
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}