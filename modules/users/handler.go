package users

import (
	"net/http"
	"stable/packages/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}



type Handler interface {
	// GetAllUsersHandler(c *gin.Context)
	GetUserByIDHandler(c *gin.Context)
	// UpdateUserHandler(c *gin.Context)
	// SetRoleHandler(c *gin.Context)
	// DeleteUserHandler(c *gin.Context)
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

func (h *handler) GetUserByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildResponseFailed("User not found", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User retrieved", user))

}
