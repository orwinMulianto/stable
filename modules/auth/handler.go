package auth

import (
	"fmt"
	"log"
	"net/http"
	"stable/packages/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

type Handler interface {
    RegisterHandler(c *gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) RegisterHandler(c *gin.Context) {
	var input RegisterInput
	var validationErrors []utils.ValidationError

	// ✅ INI YANG KURANG
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
			Field: "username",
			Error: "username cannot be empty",
		})
	}
	if input.Email == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email",
			Error: "email cannot be empty",
		})
	}
	if input.Password == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password",
			Error: "password cannot be empty",
		})
	}

	if input.Email != "" && !strings.Contains(input.Email, "@") {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "email",
			Error: "email format is invalid",
		})
	}

	if input.Password != "" && len(input.Password) < 6 {
		validationErrors = append(validationErrors, utils.ValidationError{
			Field: "password",
			Error: "password must be at least 6 characters",
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
				Field: "email",
				Error: "email is already registered",
			})
			c.JSON(http.StatusConflict, utils.BuildValidationErrorResponse("Registration failed", validationErrors))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed("Registration failed", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess("User registered successfully", UserResponseJSON(user)))
}


func VerifyEmailHandler(c *gin.Context) {
	var input VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{}
		if strings.Contains(err.Error(), "Email") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "email", Error: "email is inputuired and must be valid"})
		}
		if strings.Contains(err.Error(), "Code") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "code", Error: "verification code is inputuired"})
		}
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Verification failed", err.Error(), nil))
		return
	}
}

func ResendVerificationCodeHandler(c *gin.Context) {
    var input struct {
        Email string `json:"email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "Invalid inputuest"})
        return
    }

    code := utils.GenerateVerificationCode()

    body := fmt.Sprintf("Your verification code is: %s", code)

    log.Println("Calling SendEmail...")

    err := utils.SendEmail(input.Email, "Email Verification", body)
    if err != nil {
        log.Println("Email error:", err)
    }

    c.JSON(200, gin.H{
        "status":  true,
        "message": "Verification code sent successfully",
    })
}