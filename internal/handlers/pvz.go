package handlers

import (
	"errors"
	"net/http"
	"pvz_server/internal/app/model"
	"pvz_server/internal/app/store"

	"github.com/gin-gonic/gin"
)

type PVZInput struct {
	City model.City `json:"city" binding:"required"`
}

func CreatePVZ(storeInst store.PVZCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")

		if !ok || role != "moderator" {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		var req PVZInput

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid reques"})
			return
		}

		pvz, err := storeInst.CreatePVZ(c.Request.Context(), req.City)

		switch {
		case errors.Is(err, store.ErrCityNotAllowed):
			c.JSON(http.StatusBadRequest, gin.H{"message": "unsupported city"})
		case errors.Is(err, store.ErrDatabase):
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to create PVZ"})
		case err != nil:
			c.JSON(http.StatusBadRequest, gin.H{"message": "unexpected error"})
		default:
			c.JSON(http.StatusCreated, pvz)
		}
	}
}
