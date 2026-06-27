package trainer

import (
	"net/http"
	"strconv"

	"stable/packages/utils"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetMyProfileHandler(c *gin.Context)
	GetAllTrainersHandler(c *gin.Context)
	GetTrainerByIDHandler(c *gin.Context)
	CreateProfileHandler(c *gin.Context)
	UpdateProfileHandler(c *gin.Context)
	DeleteProfileHandler(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func getUserIDFromContext(c *gin.Context) (uint, bool) {
    value, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }

    id, ok := value.(int)
    if !ok {
        return 0, false
    }

    return uint(id), true
}
func unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, utils.BuildResponseFailed(
		"Unauthorized",
		"user not found in context",
		nil,
	))
}

func (h *handler) GetMyProfileHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		unauthorized(c)
		return
	}

	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildResponseFailed(
			"Profile not found",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess(
		"Trainer profile retrieved",
		profile,
	))
}

func (h *handler) GetAllTrainersHandler(c *gin.Context) {
	trainers, err := h.service.GetAllTrainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(
			"Failed to get trainers",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess(
		"Trainers retrieved",
		trainers,
	))
}

func (h *handler) GetTrainerByIDHandler(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(
			"Invalid ID",
			"ID must be a number",
			nil,
		))
		return
	}

	// Service baru (berdasarkan trainer ID)
	profile, err := h.service.GetTrainerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildResponseFailed(
			"Trainer not found",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess(
		"Trainer retrieved",
		profile,
	))
}

func (h *handler) CreateProfileHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		unauthorized(c)
		return
	}

	var req CreateTrainerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(
			"Invalid request",
			err.Error(),
			nil,
		))
		return
	}

	profile, err := h.service.CreateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(
			"Failed to create profile",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildResponseSuccess(
		"Trainer profile created",
		profile,
	))
}

func (h *handler) UpdateProfileHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		unauthorized(c)
		return
	}

	var req UpdateTrainerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(
			"Invalid request",
			err.Error(),
			nil,
		))
		return
	}

	profile, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed(
			"Failed to update profile",
			err.Error(),
			nil,
		))
		return
	}
	
	c.JSON(http.StatusOK, utils.BuildResponseSuccess(
		"Trainer profile updated",
		profile,
	))
}

func (h *handler) DeleteProfileHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		unauthorized(c)
		return
	}

	if err := h.service.DeleteProfile(userID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(
			"Failed to delete profile",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess(
		"Trainer profile deleted",
		nil,
	))
}