package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"pvz_server/internal/app/apiserver"
	"pvz_server/internal/app/deps"
	"pvz_server/internal/app/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullReceptionFlow(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	s := apiserver.NewServerWithDeps(&deps.Dependencies{
		Store: store.New(db),
	})

	ts := httptest.NewServer(s.GetEngine())
	defer ts.Close()

	employeeToken := getToken(t, ts.URL, "employee")
	moderatorToken := getToken(t, ts.URL, "moderator")

	pvzID := createPVZ(t, ts.URL, moderatorToken, "Казань")

	createReception(t, ts.URL, employeeToken, pvzID)

	for i := 0; i < 20; i++ {
		var productType string
		switch {
		case i < 20:
			productType = "электроника"
		case i < 40:
			productType = "одежда"
		default:
			productType = "обувь"
		}
		addProduct(t, ts.URL, employeeToken, pvzID, productType)
	}

	closeReception(t, ts.URL, employeeToken, pvzID)
}

func getToken(t *testing.T, baseURL, role string) string {
	body := map[string]string{"role": role}
	data, _ := json.Marshal(body)

	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(data))

	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}

	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	token := result["token"]
	assert.NotEmpty(t, token)
	return token
}

func createPVZ(t *testing.T, baseURL, token, city string) string {
	body := map[string]string{"city": city}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/pvz", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("failed to create PVZ: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	id, ok := result["id"].(string)

	if !ok || id == "" {
		t.Fatal("invalid pvz ID")
	}

	return id
}

func createReception(t *testing.T, baseURL, token, pvzID string) {
	body := map[string]string{"pvzId": pvzID}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/receptions", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("failed to create reception: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func addProduct(t *testing.T, baseURL, token, pvzID, productType string) {
	body := map[string]string{"pvzId": pvzID, "type": productType}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/products", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("failed to add product: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func closeReception(t *testing.T, baseURL, token, pvzID string) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("failed to close reception: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
