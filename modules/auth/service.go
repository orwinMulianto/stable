package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"stable/database/entities"
	"stable/modules/users"
	"stable/packages/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service interface {
	Register(input RegisterInput) (entities.User, error)
	Login(input LoginInput) (string, entities.User, error)
	VerifyEmail(email, code string) error
	ResendVerificationCode(email string) error
    GetGoogleAuthURL() string                          // ← tambah
    GoogleLogin(code string) (string, entities.User, error) 
}

type authService struct {
	repo        users.Repository
	oauthConfig *oauth2.Config
}


func NewService(repo users.Repository) Service {
	return &authService{
        repo: repo,
        oauthConfig: &oauth2.Config{
            ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
            ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
            RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
            Scopes: []string{
                "https://www.googleapis.com/auth/userinfo.email",
                "https://www.googleapis.com/auth/userinfo.profile",
            },
            Endpoint: google.Endpoint,
        },
    }
}


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

	go func() {
		emailBody := "Your verification code is: " + code
		if err := utils.SendEmail(user.Email, "Email Verification", emailBody); err != nil {
			log.Println("[EMAIL] FAILED:", err)
		}
	}()

	return user, nil
}

func (s *authService) Login(input LoginInput) (string, entities.User, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", entities.User{}, errors.New("email not found")
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", entities.User{}, errors.New("invalid password")
	}

	if !user.IsVerified {
		return "", user, errors.New("email not verified")
	}

	// Generate JWT token
	token, err := generateToken(int(user.ID))
	if err != nil {
		return "", entities.User{}, err
	}

	return token, user, nil
}

func (s *authService) ResendVerificationCode(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	code := utils.GenerateVerificationCode()
	log.Println("[RESEND] Plain code:", code) // ← tambah ini

	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedCodeStr := string(hashedCode)
	user.VerificationCode = &hashedCodeStr
	log.Println("[RESEND] Hashed code:", hashedCodeStr) // ← tambah ini

	if err := s.repo.Update(user); err != nil {
		log.Println("[RESEND] Update error:", err) // ← tambah ini
		return err
	}

	log.Println("[RESEND] Update success") // ← tambah ini

	// Cek dari DB setelah update
	updatedUser, _ := s.repo.FindByEmail(email)
	log.Println("[RESEND] Code in DB after update:", updatedUser.VerificationCode) // ← tambah ini

	emailBody := `<h2>Email Verification</h2>
        <p>Your verification code is: <strong>` + code + `</strong></p>`
	go utils.SendEmail(user.Email, "Email Verification", emailBody)

	return nil
}

func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) VerifyEmail(email, code string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	if user.VerificationCode == nil {
		return errors.New("invalid verification code")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.VerificationCode), []byte(code))
	if err != nil {
		return errors.New("invalid verification code")
	}

	user.IsVerified = true
	user.VerificationCode = nil

	return s.repo.Update(user)
}

func (s *authService) GetGoogleAuthURL() string {
    return s.oauthConfig.AuthCodeURL("state-token")
}

func (s *authService) GoogleLogin(code string) (string, entities.User, error) {
    oauthToken, err := s.oauthConfig.Exchange(context.Background(), code)
    if err != nil {
        return "", entities.User{}, errors.New("failed to exchange token: " + err.Error())
    }

    client := s.oauthConfig.Client(context.Background(), oauthToken)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        return "", entities.User{}, errors.New("failed to get user info: " + err.Error())
    }
    defer resp.Body.Close()

    var googleUser GoogleUserInfo
    if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
        return "", entities.User{}, errors.New("failed to decode user info: " + err.Error())
    }

    user, err := s.repo.FindByEmail(googleUser.Email)
    if err != nil {
        newUser := &entities.User{  
            Email:      googleUser.Email,
            Username:   googleUser.Name,
            IsVerified: true,
        }

        if err := s.repo.Create(newUser); err != nil {
            return "", entities.User{}, errors.New("failed to create user: " + err.Error())
        }
        user = *newUser
    }

    token, err := utils.GenerateToken(user)
	if err != nil {
    	return "", entities.User{}, errors.New("failed to generate token: " + err.Error())
	}

    return token, user, nil
}