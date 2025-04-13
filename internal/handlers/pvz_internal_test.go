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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	createFunc func(ctx context.Context, city model.City) (*model.PVZ, error)
}

func (m *mockStore) CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error) {
	return m.createFunc(ctx, city)
}

func setupRouterWithRole(role string, store storeInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})

	r.POST("/pvz", handlers.CreatePVZ(store.(*mockStore)))
	return r
}

type storeInterface interface {
	CreatePVZ(ctx context.Context, city model.City) (*model.PVZ, error)
}

func TestCreatePVZ_Success(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, city model.City) (*model.PVZ, error) {
			return &model.PVZ{City: city}, nil
		},
	}

	router := setupRouterWithRole("moderator", mock)

	body := map[string]string{"city": "Москва"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Москва")
}

func TestCreatePVZ_NoRole(t *testing.T) {
	mock := &mockStore{}
	router := setupRouterWithRole("", mock)

	req, _ := http.NewRequest("POST", "/pvz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "access denied")
}

func TestCreatePVZ_InvalidRole(t *testing.T) {
	mock := &mockStore{}
	router := setupRouterWithRole("employee", mock)

	req, _ := http.NewRequest("POST", "/pvz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "access denied")
}

func TestCreatePVZ_InvalidBody(t *testing.T) {
	mock := &mockStore{}
	router := setupRouterWithRole("moderator", mock)

	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid")
}

func TestCreatePVZ_UnsupportedCity(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, city model.City) (*model.PVZ, error) {
			switch city {
			case "Москва", "Санкт-Петербург", "Казань":
				return &model.PVZ{City: city}, nil
			default:
				return nil, store.ErrCityNotAllowed
			}
		},
	}
	router := setupRouterWithRole("moderator", mock)

	body := map[string]string{"city": "Сочи"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "unsupported city")
}

func TestCreatePVZ_DatabaseError(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, city model.City) (*model.PVZ, error) {
			return nil, store.ErrDatabase
		},
	}
	router := setupRouterWithRole("moderator", mock)

	body := map[string]string{"city": "Казань"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create")
}

func TestCreatePVZ_UnexpectedError(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, city model.City) (*model.PVZ, error) {
			return nil, errors.New("something went wrong")
		},
	}
	router := setupRouterWithRole("moderator", mock)

	body := map[string]string{"city": "Казань"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/pvz", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "unexpected error")
}
