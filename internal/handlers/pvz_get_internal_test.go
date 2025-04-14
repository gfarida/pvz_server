package handlers_test

import (
	"context"
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

type mockPVZFetcher struct {
	fetchFunc func(ctx context.Context, start, end *time.Time, page, limit int) ([]*model.PVZWithReceptions, error)
}

func (m *mockPVZFetcher) FetchPVZList(ctx context.Context, start, end *time.Time, page, limit int) ([]*model.PVZWithReceptions, error) {
	return m.fetchFunc(ctx, start, end, page, limit)
}

func setupPVZGetRouter(role string, fetcher *mockPVZFetcher) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})

	r.GET("/pvz", handlers.GetPVZList(fetcher))
	return r
}

func TestGetPVZList_Success(t *testing.T) {
	mock := &mockPVZFetcher{
		fetchFunc: func(ctx context.Context, start, end *time.Time, page, limit int) ([]*model.PVZWithReceptions, error) {
			return []*model.PVZWithReceptions{
				{
					PVZ: model.PVZ{City: "Москва"},
				},
			}, nil
		},
	}

	router := setupPVZGetRouter("employee", mock)
	req, _ := http.NewRequest("GET", "/pvz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Москва")
}

func TestGetPVZList_InvalidRole(t *testing.T) {
	mock := &mockPVZFetcher{}
	router := setupPVZGetRouter("guest", mock)

	req, _ := http.NewRequest("GET", "/pvz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "access denied")
}

func TestGetPVZList_InvalidDate(t *testing.T) {
	mock := &mockPVZFetcher{}
	router := setupPVZGetRouter("moderator", mock)

	req, _ := http.NewRequest("GET", "/pvz?startDate=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid startDate")
}

func TestGetPVZList_InvalidPagination(t *testing.T) {
	mock := &mockPVZFetcher{}
	router := setupPVZGetRouter("employee", mock)

	req, _ := http.NewRequest("GET", "/pvz?page=0&limit=1000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid pagination")
}

func TestGetPVZList_DatabaseError(t *testing.T) {
	mock := &mockPVZFetcher{
		fetchFunc: func(ctx context.Context, start, end *time.Time, page, limit int) ([]*model.PVZWithReceptions, error) {
			return nil, store.ErrDatabase
		},
	}
	router := setupPVZGetRouter("moderator", mock)

	req, _ := http.NewRequest("GET", "/pvz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to fetch PVZ list")
}
