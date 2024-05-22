package test

import (
	"bytes"
	sqlc "client_task/pkg/common/db/sqlc"
	"client_task/pkg/controllers"
	"client_task/pkg/payloads"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseUserPath string = "/api/users"

func TestCreateUser(t *testing.T) {
	router := setupTestRouter()

	payload := payloads.CreateUser{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "password123",
		UserType:     "jobseeker",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseUserPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfully created user", response["status"])

	// Clean up
	truncateTables(t, "users")
}

func TestUpdateUser(t *testing.T) {
	router := setupTestRouter()

	// Create a user to update
	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "jobseeker",
	})
	assert.NoError(t, err)

	userID := strconv.Itoa(int(user.ID))

	payload := payloads.UpdateUser{
		FullName:     "Jane Doe",
		PasswordHash: "newpassword123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPut, baseUserPath+"/"+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfully updated User", response["status"])

	// Clean up
	truncateTables(t, "users")
}

func TestGetUserById(t *testing.T) {
	router := setupTestRouter()

	// Create a user to retrieve
	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	userID := strconv.Itoa(int(user.ID))

	req, _ := http.NewRequest(http.MethodGet, baseUserPath+"/"+userID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrived id", response["status"])

	// Clean up
	truncateTables(t, "users")
}

func TestGetAllUsers(t *testing.T) {
	router := setupTestRouter()

	// Create a user to ensure there's at least one user in the database
	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	_, err = testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, baseUserPath, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrieved all Users", response["status"])
	assert.Greater(t, int(response["size"].(float64)), 0)

	// Clean up
	truncateTables(t, "users")
}

func TestDeleteUserById(t *testing.T) {
	router := setupTestRouter()

	// Create a user to delete
	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	userID := strconv.Itoa(int(user.ID))

	req, _ := http.NewRequest(http.MethodDelete, baseUserPath+"/"+userID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfuly deleted", response["status"])

	// Clean up
	truncateTables(t, "users")
}
