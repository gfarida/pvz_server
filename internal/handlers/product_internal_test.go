package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pvz_server/internal/app/model"
	"pvz_server/internal/app/store"
	"pvz_server/internal/handlers"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockProductStore struct {
	addFunc func(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error)
}

func (m *mockProductStore) AddProduct(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error) {
	return m.addFunc(ctx, pvzID, productType)
}

func setupProductRouterWithRole(role string, store *mockProductStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})

	r.POST("/products", handlers.AddProduct(store))
	return r
}

func TestAddProduct_Success(t *testing.T) {
	mock := &mockProductStore{
		addFunc: func(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error) {
			return &model.Product{
				ID:          "p-123",
				DateTime:    time.Now(),
				Type:        productType,
				ReceptionID: "r-123",
			}, nil
		},
	}

	router := setupProductRouterWithRole("employee", mock)

	body := map[string]string{"type": "электроника", "pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"type":"электроника"`)
}

func TestAddProduct_InvalidRole(t *testing.T) {
	mock := &mockProductStore{}
	router := setupProductRouterWithRole("moderator", mock)

	body := map[string]string{"type": "электроника", "pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "access denied")
}

func TestAddProduct_InvalidBody(t *testing.T) {
	mock := &mockProductStore{}
	router := setupProductRouterWithRole("employee", mock)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestAddProduct_NoActiveReception(t *testing.T) {
	mock := &mockProductStore{
		addFunc: func(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error) {
			return nil, store.ErrNoActiveReception
		},
	}
	router := setupProductRouterWithRole("employee", mock)

	body := map[string]string{"type": "электроника", "pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "no active reception")
}

func TestAddProduct_DatabaseError(t *testing.T) {
	mock := &mockProductStore{
		addFunc: func(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error) {
			return nil, store.ErrDatabase
		},
	}
	router := setupProductRouterWithRole("employee", mock)

	body := map[string]string{"type": "электроника", "pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to add product")
}

func TestAddProduct_UnexpectedError(t *testing.T) {
	mock := &mockProductStore{
		addFunc: func(ctx context.Context, pvzID string, productType model.ProductType) (*model.Product, error) {
			return nil, errors.New("unknown")
		},
	}
	router := setupProductRouterWithRole("employee", mock)

	body := map[string]string{"type": "электроника", "pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "unexpected error")
}
