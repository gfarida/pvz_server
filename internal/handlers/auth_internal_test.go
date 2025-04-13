package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"pvz_server/internal/handlers"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDummyLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	router := gin.Default()
	router.POST("/dummyLogin", handlers.DummyLogin)

	payload := map[string]string{"role": "moderator"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/dummyLogin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &result)

	assert.NoError(t, err)
	assert.Contains(t, result, "token")
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	router := gin.Default()
	router.POST("/dummyLogin", handlers.DummyLogin)

	payload := map[string]string{"role": "test"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/dummyLogin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestDummyLogin_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	router := gin.Default()
	router.POST("/dummyLogin", handlers.DummyLogin)

	req, _ := http.NewRequest("POST", "/dummyLogin", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
