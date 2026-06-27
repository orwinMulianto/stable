package dailychallenge

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTodayHandler(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "valid user_id is required"})
		return
	}

	response, err := h.service.GetToday(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to load daily challenge",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) CompleteTodayHandler(c *gin.Context) {
	var request CompleteChallengeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.service.CompleteToday(request.UserID, request.Repetitions)
	if errors.Is(err, ErrAlreadyCompleted) {
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}
	if errors.Is(err, ErrTargetNotReached) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to complete daily challenge",
			"error":   err.Error(),
	})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}
