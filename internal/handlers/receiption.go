package handlers

import (
	"errors"
	"net/http"
	"pvz_server/internal/app/store"

	"github.com/gin-gonic/gin"
)

type ReceprionInput struct {
	PVZID string `json:"pvzId" binding:"required"`
}

func CreateReception(storeInst store.ReceptionCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")

		if !ok || role != "employee" {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		var req ReceprionInput

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		reception, err := storeInst.CreateReception(c.Request.Context(), req.PVZID)

		switch {
		case err == store.ErrReceptionAlreadyExists:
			c.JSON(http.StatusBadRequest, gin.H{"message": "previous reception is not closed"})
		case err == store.ErrDatabase:
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to create reception"})
		case err != nil:
			c.JSON(http.StatusBadRequest, gin.H{"message": "unexpected error"})
		default:
			c.JSON(http.StatusCreated, reception)
		}
	}
}

func CloseLastReception(storeInst store.ReceptionCloser) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")

		if !ok || role != "employee" {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		pvzID := c.Param("pvzId")

		if pvzID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "pvzId is required"})
			return
		}

		reception, err := storeInst.CloseLastReception(c.Request.Context(), pvzID)

		switch {
		case errors.Is(err, store.ErrNoActiveReception):
			c.JSON(http.StatusBadRequest, gin.H{"message": "no active reception to close"})
		case errors.Is(err, store.ErrDatabase):
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to close reception"})
		case err != nil:
			c.JSON(http.StatusBadRequest, gin.H{"message": "unexpected error"})
		default:
			c.JSON(http.StatusOK, reception)
		}
	}
}
