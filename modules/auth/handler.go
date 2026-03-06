package auth

import (
	"net/http"
	"stable/database/entities"
	"stable/database/migrations"
	"stable/packages/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	service Service
}

func RegisterHandler(c *gin.Context) {
    var req RegisterInput

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

    user := entities.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
    }

    db := migrations.GetDB()

    result := db.Create(&user)

    if result.Error != nil {
        c.JSON(500, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(200, gin.H{
        "message": "Akun berhasil dibuat",
        "user_id": user.ID,
    })
}


func VerifyEmailHandler(c *gin.Context) {
	var input VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		validationErrors := []utils.ValidationError{}
		if strings.Contains(err.Error(), "Email") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "email", Error: "email is required and must be valid"})
		}
		if strings.Contains(err.Error(), "Code") {
			validationErrors = append(validationErrors, utils.ValidationError{Field: "code", Error: "verification code is required"})
		}
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, utils.BuildValidationErrorResponse("Validation failed", validationErrors))
			return
		}
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Verification failed", err.Error(), nil))
		return
	}
}
	