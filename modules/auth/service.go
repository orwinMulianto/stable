package auth

import (
	"stable/database/entities"
	"stable/packages/utils"

	"golang.org/x/crypto/bcrypt"
	"log"
	// "google.golang.org/api/oauth2/v2"
	"errors"
	"stable/modules/users"
)


type Service interface {
	Register(input RegisterInput) (entities.User, error)
	// Login(input LoginInput) (string, error)
	// VerifyEmail(email, code string) error
	ResendVerificationCode(email string) error
	// Me(user_id int) (*entities.User, error)
	// LoginOrRegisterWithGoogle(googleUserInfo *oauth2.Userinfo) (*entities.User, string, error)
	// ForgotPassword(input ForgotPasswordInput) error
	// ResetPassword(input ResetPasswordInput) error
}

type authService struct {
	repo users.Repository
}

func NewService(repo users.Repository) Service {
	return &authService{repo: repo}
}

func (s *authService) Register(input RegisterInput) (entities.User, error) {
	defaultRole := "user"
	user := entities.User{
		Username:     input.Username,
		Email:        input.Email,
		AuthProvider: "Password",
		Role:         defaultRole,
		IsVerified: false,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(hashedPassword)

	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return user, errors.New("email already registered")
	}
	code := utils.GenerateVerificationCode()
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	hashedCodeStr := string(hashedCode)
	user.VerificationCode = &hashedCodeStr

	err = s.repo.Create(&user)
	if err != nil {
		return user, err
	}

	emailBody := "Your verification code is: " + code
	go func() {
		err := utils.SendEmail(user.Email, "Email Verification", emailBody)
		if err != nil {
			log.Println("[EMAIL] FAILED:", err)
		}
	}()

	return user, nil
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
	go utils.SendEmail(user.Email, "Email Verification - Reduka", emailBody)

	return nil
}
