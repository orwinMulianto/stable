package trainerchat

import (
	"errors"
	"net/http"
	"strconv"
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

func getUserIDFromContext(c *gin.Context) (uint, bool) {
    value, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }

    switch v := value.(type) {
    case int:
        if v <= 0 {
            return 0, false
        }
        return uint(v), true
    case uint:
        return v, true
    case float64:
        return uint(v), true
    default:
        return 0, false
    }
}

func (h *Handler) ListTrainersHandler(c *gin.Context) {
	trainers, err := h.service.ListTrainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to load trainers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trainers})
}

func (h *Handler) DashboardMeHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	response, err := h.service.GetMyDashboard(userID)
	if errors.Is(err, ErrTrainerNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "trainer not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to load trainer dashboard",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) CheckoutHandler(c *gin.Context) {
	var request CheckoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid checkout request", "error": err.Error()})
		return
	}

	response, err := h.service.Checkout(request)
	if errors.Is(err, ErrTrainerNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "trainer not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create payment", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) NotificationHandler(c *gin.Context) {
	var notification MidtransNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid notification", "error": err.Error()})
		return
	}

	if err := h.service.HandleNotification(notification); err != nil {
		if errors.Is(err, ErrInvalidNotification) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid notification signature"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to process notification", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification processed"})
}

func (h *Handler) HistoryHandler(c *gin.Context) {
    userID, ok := getUserIDFromContext(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }

    response, err := h.service.GetHistory(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load history", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) SessionHandler(c *gin.Context) {
	sessionID, err := parseUintParam(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
		return
	}

	response, err := h.service.GetSession(sessionID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "session not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load session", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) SendMessageHandler(c *gin.Context) {
	sessionID, err := parseUintParam(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
		return
	}

	var request SendMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid message request", "error": err.Error()})
		return
	}

	response, err := h.service.SendMessage(sessionID, request)
	if errors.Is(err, ErrSessionNotAvailable) {
		c.JSON(http.StatusForbidden, gin.H{"message": "chat session is not active"})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "session not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to send message", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func parseUintParam(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	if parsed == 0 {
		return 0, errors.New("id must be greater than zero")
	}

	return uint(parsed), nil
}

func (h *Handler) GetTrainerSessionsHandler(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	response, err := h.service.GetTrainerSessions(userID)
	if errors.Is(err, ErrTrainerNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "trainer not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to load trainer sessions",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) ConfirmPaymentHandler(c *gin.Context) {
    sessionID, err := parseUintParam(c.Param("session_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
        return
    }

    now := time.Now()
    startedAt := now
    expiresAt := now.Add(10 * time.Minute)

    updates := map[string]interface{}{
        "status":     "paid",
        "paid_at":    now,
        "started_at": startedAt,
        "expires_at": expiresAt,
    }

    if err := h.service.UpdatePaymentDirect(sessionID, updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to confirm payment"})
        return
    }

    session, _ := h.service.GetSession(sessionID)
    c.JSON(http.StatusOK, gin.H{"data": session})
}

func (h *Handler) DevMarkPaidHandler(c *gin.Context) {
	sessionID, err := parseUintParam(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
		return
	}

	response, err := h.service.DevMarkPaid(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to mark session paid",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}