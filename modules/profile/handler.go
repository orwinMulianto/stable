package profile

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfileHandler(c *gin.Context) {
	userID, err := parseUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "valid user_id is required"})
		return
	}

	response, err := h.service.GetProfile(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to load profile",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) UpdateProfileHandler(c *gin.Context) {
	userID, err := parseUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "valid user_id is required"})
		return
	}

	var request UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.service.UpdateProfile(userID, request)
	var cooldownErr UsernameCooldownError
	if errors.As(err, &cooldownErr) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message":                "username can only be changed once every 14 days",
			"can_change_username_at": cooldownErr.CanChangeAt,
		})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) UploadAvatarHandler(c *gin.Context) {
	userID, err := parseUserID(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "valid user_id is required"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "avatar file is required"})
		return
	}

	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "avatar max size is 2MB"})
		return
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExtensions[extension] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "avatar must be jpg, png, or webp"})
		return
	}

	uploadDir := filepath.Join("uploads", "profile")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to prepare upload folder",
			"error":   err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("%d-%d%s", userID, time.Now().UnixNano(), extension)
	storagePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, storagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save avatar",
			"error":   err.Error(),
		})
		return
	}

	publicPath := "/" + filepath.ToSlash(storagePath)
	response, err := h.service.UpdateProfileImage(userID, publicPath)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update avatar",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       response,
		"avatar_url": publicPath,
	})
}

func parseUserID(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	if parsed == 0 {
		return 0, errors.New("user_id must be greater than zero")
	}

	return uint(parsed), nil
}
