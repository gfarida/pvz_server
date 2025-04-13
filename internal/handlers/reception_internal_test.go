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

type mockReceptionStore struct {
	createFunc func(ctx context.Context, pvzID string) (*model.Reception, error)
}

func (m *mockReceptionStore) CreateReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	return m.createFunc(ctx, pvzID)
}

type receptionStoreInterface interface {
	CreateReception(ctx context.Context, pvzID string) (*model.Reception, error)
}

func setupReceptionRouterWithRole(role string, store receptionStoreInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})

	r.POST("/receptions", handlers.CreateReception(store.(*mockReceptionStore)))
	return r
}

func TestCreateReception_Success(t *testing.T) {
	mock := &mockReceptionStore{
		createFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
			return &model.Reception{
				ID:       "rec-123",
				PvzID:    pvzID,
				DateTime: time.Now(),
				Status:   model.InProgress,
			}, nil
		},
	}
	router := setupReceptionRouterWithRole("employee", mock)

	body := map[string]string{"pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"pvzId":"pvz1"`)
}

func TestCreateReception_InvalidRole(t *testing.T) {
	mock := &mockReceptionStore{}
	router := setupReceptionRouterWithRole("moderator", mock)

	body := map[string]string{"pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "access denied")
}

func TestCreateReception_InvalidBody(t *testing.T) {
	mock := &mockReceptionStore{}
	router := setupReceptionRouterWithRole("employee", mock)

	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestCreateReception_AlreadyExists(t *testing.T) {
	mock := &mockReceptionStore{
		createFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
			return nil, store.ErrReceptionAlreadyExists
		},
	}
	router := setupReceptionRouterWithRole("employee", mock)

	body := map[string]string{"pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "previous reception is not closed")
}

func TestCreateReception_DatabaseError(t *testing.T) {
	mock := &mockReceptionStore{
		createFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
			return nil, store.ErrDatabase
		},
	}
	router := setupReceptionRouterWithRole("employee", mock)

	body := map[string]string{"pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create reception")
}

func TestCreateReception_UnexpectedError(t *testing.T) {
	mock := &mockReceptionStore{
		createFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
			return nil, errors.New("something went wrong")
		},
	}
	router := setupReceptionRouterWithRole("employee", mock)

	body := map[string]string{"pvzId": "pvz1"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/receptions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "unexpected error")
}
