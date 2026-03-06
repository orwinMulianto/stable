package auth

import (
	"stable/database/entities"

	"google.golang.org/api/oauth2/v2"
)

type Service interface {
	Register(input RegisterInput) (entities.User, error)
	Login(input LoginInput) (string, error)
	VerifyEmail(email, code string) error
	ResendVerificationCode(email string) error
	Me(user_id int) (*entities.User, error)
	LoginOrRegisterWithGoogle(googleUserInfo *oauth2.Userinfo) (*entities.User, string, error)
	ForgotPassword(input ForgotPasswordInput) error
	ResetPassword(input ResetPasswordInput) error
}