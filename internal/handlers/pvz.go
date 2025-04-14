package handlers

import (
	"errors"
	"net/http"
	"pvz_server/internal/app/model"
	"pvz_server/internal/app/store"
	"strconv"
	"time"

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

func GetPVZList(storeInst store.PVZFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")

		if !ok || (role != "moderator" && role != "employee") {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		startDateStr := c.Query("startDate")
		endDateStr := c.Query("endDate")

		var startDate, endDate *time.Time

		if startDateStr != "" {
			t, err := time.Parse(time.RFC3339, startDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "invalid startDate"})
				return
			}
			startDate = &t
		}

		if endDateStr != "" {
			t, err := time.Parse(time.RFC3339, endDateStr)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "invalid endDate"})
				return
			}

			endDate = &t
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

		if page < 1 || limit < 1 || limit > 30 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid pagination"})
			return
		}

		pvzs, err := storeInst.FetchPVZList(c.Request.Context(), startDate, endDate, page, limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to fetch PVZ list"})
			return
		}

		c.JSON(http.StatusOK, pvzs)
	}
}
