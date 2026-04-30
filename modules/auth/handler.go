package auth

import (
	"fmt"
	// "log"
	"net/http"
	"stable/packages/utils"
	"strings"
	// "stable/database/migrations"
	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

type Handler interface {
    RegisterHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	VerifyEmailHandler(c *gin.Context)
	ResendVerificationCodeHandler(c *gin.Context)
	GoogleLoginHandler(c *gin.Context)       // ← tambah
    GoogleCallbackHandler(c *gin.Context)    // ← tambah
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}


func (h *handler) RegisterHandler(c *gin.Context) {
	var input RegisterInput
	var validationErrors []utils.ValidationError

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if input.Username == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "username", Error: "username cannot be empty",
		})
	}
	if input.Email == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email", Error: "email cannot be empty",
		})
	}
	if input.Password == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password", Error: "password cannot be empty",
		})
	}
	if input.Email != "" && !strings.Contains(input.Email, "@") {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email", Error: "email format is invalid",
		})
	}
	if input.Password != "" && len(input.Password) < 6 {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password", Error: "password must be at least 6 characters",
		})
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
		if strings.Contains(err.Error(), "email already registered") {
			validationErrors = append(validationErrors, utils.ValidationError{
				Field: "email", Error: "email is already registered",
			})
			c.JSON(http.StatusConflict, utils.BuildValidationErrorResponse("Registration failed", validationErrors))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Registration failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("User registered successfully", UserResponseJSON(user)))
}

func (h *handler) LoginHandler(c *gin.Context) {
	var input LoginInput
	var validationErrors []utils.ValidationError

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if input.Email == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email", Error: "email cannot be empty",
		})
	}
	if input.Password == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password", Error: "password cannot be empty",
		})
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
		return
	}

	token, user, err := h.service.Login(input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("Login failed", "Email tidak ditemukan", nil))
			return
		}
		if strings.Contains(err.Error(), "not verified") {
			c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("Login failed", "Email belum diverifikasi", nil))
			return
		}
		if strings.Contains(err.Error(), "invalid password") {
			c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed("Login failed", "Password salah", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Login failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("Login berhasil", gin.H{
		"token": token,
		"user":  UserResponseJSON(user),
	}))
}


func (h *handler) ResendVerificationCodeHandler(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    if err := h.service.ResendVerificationCode(input.Email); err != nil {
        c.JSON(400, gin.H{
            "status":  false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "status":  true,
        "message": "Verification code sent successfully",
    })
}

func (h *handler) VerifyEmailHandler(c *gin.Context) {
	var input VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{}

		if strings.Contains(err.Error(), "Email") {
			validationErrors = append(validationErrors, utils.ValidationError{
				Field: "email", Error: "email is required and must be valid",
			})
		}
		if strings.Contains(err.Error(), "Code") {
			validationErrors = append(validationErrors, utils.ValidationError{
				Field: "code", Error: "verification code is required",
			})
		}

		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest,
				utils.BuildValidationErrorResponse("Validation failed", validationErrors))
			return
		}

		c.JSON(http.StatusBadRequest,
			utils.BuildResponseFailed("Verification failed", err.Error(), nil))
		return
	}
	err := h.service.VerifyEmail(input.Email, input.Code)
	if err != nil {

		switch err.Error() {

		case "email not found":
			c.JSON(http.StatusNotFound,
				utils.BuildValidationErrorResponse("Email not found", []utils.ValidationError{
					{Field: "email", Error: "email is not registered"},
				}))
			return

		case "email already verified":
			c.JSON(http.StatusBadRequest,
				utils.BuildResponseFailed("Already verified", "email is already verified", nil))
			return

		case "invalid verification code":
			c.JSON(http.StatusBadRequest,
				utils.BuildValidationErrorResponse("Verification failed", []utils.ValidationError{
					{Field: "code", Error: "invalid verification code"},
				}))
			return
		}
		c.JSON(http.StatusBadRequest,
			utils.BuildResponseFailed("Verification failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK,
		utils.BuildResponseSuccess("Email verified successfully", nil))
}

func (h *handler) GoogleLoginHandler(c *gin.Context) {
    url := h.service.GetGoogleAuthURL()
    c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *handler) GoogleCallbackHandler(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(
            "Google login failed", "authorization code not found", nil,
        ))
        return
    }

    token, user, err := h.service.GoogleLogin(code)
    if err != nil {
        c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5501/login?error=google_failed")
        return
    }

    _ = user

    redirectURL := fmt.Sprintf("http://localhost:5501/page/dashboard.html?token=%s", token)
    c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}