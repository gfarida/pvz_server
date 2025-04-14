package handlers

import (
	"net/http"
	"pvz_server/internal/app/store"

	"github.com/gin-gonic/gin"
)

type ReceprionInput struct {
	PVZID string `json:"pvzId" binding:"required"`
}

func CreateReception(storeInst store.ReceptionCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role != "employee" {
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
