package user

import (
	"net/http"
	"strconv"
	"stable/packages/utils"
	"github.com/gin-gonic/gin"
)

func GetUserbyIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildResponseFailed("Invalid ID", "ID must be a number", nil))
		return
	}

	user, err := GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildResponseFailed("User not found", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseSuccess("User retrieved", user))





}