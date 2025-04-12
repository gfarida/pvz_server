package handlers

import (
	"net/http"
	"pvz_server/internal/utils"

	"github.com/gin-gonic/gin"
)

type dummyLoginRequest struct {
	Role string `json:"role" binding:"required,oneof=employee moderator"`
}

func DummyLogin(c *gin.Context) {
	var req dummyLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	token, err := utils.GenerateJWT(req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
