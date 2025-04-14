package handlers

import (
	"errors"
	"net/http"
	"pvz_server/internal/app/model"
	"pvz_server/internal/app/store"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Type  model.ProductType `json:"type" binding:"required"`
	PVZID string            `json:"pvzId" binding:"required"`
}

func AddProduct(storeInst store.ProductAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")

		if !ok || role != "employee" {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		var req ProductInput

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		product, err := storeInst.AddProduct(c.Request.Context(), req.PVZID, req.Type)

		switch {
		case errors.Is(err, store.ErrProductTypeNotAllowed):
			c.JSON(http.StatusBadRequest, gin.H{"message": "unsupported product type"})
		case errors.Is(err, store.ErrNoActiveReception):
			c.JSON(http.StatusBadRequest, gin.H{"message": "no active reception"})
		case errors.Is(err, store.ErrDatabase):
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to add product"})
		case err != nil:
			c.JSON(http.StatusBadRequest, gin.H{"message": "unexpected error"})
		default:
			c.JSON(http.StatusCreated, product)
		}
	}
}
