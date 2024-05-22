package test

import (
	"bytes"
	sqlc "client_task/pkg/common/db/sqlc"
	"client_task/pkg/controllers"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var base_path string = "/api/auth"

func TestLoginHandler(t *testing.T) {
	router := setupTestRouter()

	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)

	// Insert a user into the test database
	_, err = testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "Test User",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		UserType:     "jobseeker",
	})
	assert.NoError(t, err)

	// Valid login
	loginReq := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, base_path+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Login successfully", response["status"])

	// Invalid login
	loginReq["password"] = "wrongpassword"
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest(http.MethodPost, base_path+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid credentials", response["error"])

	truncateTables(t, "users")
}

func TestRegisterHandler(t *testing.T) {
	router := setupTestRouter()

	// Valid registration
	registerReq := map[string]string{
		"full_name": "John Doe",
		"email":     "john@example.com",
		"password":  "password123",
		"user_type": "user",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, base_path+"/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Registration successful", response["status"])
	assert.NotNil(t, response["user"])

	// Invalid registration payload
	registerReq["email"] = ""
	body, _ = json.Marshal(registerReq)
	req, _ = http.NewRequest(http.MethodPost, base_path+"/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["error"])
}
