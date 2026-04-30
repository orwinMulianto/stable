package auth

import (
	"stable/database/entities"
)

func UserResponseJSON(user entities.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		NoTelp:       user.NoTelp,
		JenisKelamin: user.JenisKelamin,
		AuthProvider: user.AuthProvider,
		ProfileImage: user.ProfileImage,
		IsVerified:   user.IsVerified,
	}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type ResendVerificationInput struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UserResponse struct {
	ID           uint    `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	NoTelp       string  `json:"noTelp,omitempty"`
	JenisKelamin *bool   `json:"jenisKelamin,omitempty"`
	AuthProvider string  `json:"authProvider"`
	ProfileImage string  `json:"profileImage,omitempty"`
	IsVerified   bool    `json:"isVerified"`
}

type GoogleUserInfo struct {
    ID      string `json:"id"`
    Email   string `json:"email"`
    Name    string `json:"name"`
    Picture string `json:"picture"`
}