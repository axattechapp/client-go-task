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

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

var baseProfilePath string = "/api/profiles"

func TestCreateProfile(t *testing.T) {
	router := setupTestRouter()

	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	payload := payloads.CreateProfile{
		UserID:      user.ID,
		Bio:         pgtype.Text{String: "This is a bio", Valid: true},
		Company:     "Example Company",
		JobRole:     "Developer",
		Description: pgtype.Text{String: "This is a description", Valid: true},
		Skills:      []string{"Go", "Gin"},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseProfilePath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfully created profile", response["status"])

	// Clean up
	truncateTables(t, "profiles", "users")
}

func TestUpdateProfile(t *testing.T) {
	router := setupTestRouter()

	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	profile, err := testQueries.CreateProfile(context.TODO(), sqlc.CreateProfileParams{
		UserID:      user.ID,
		Bio:         pgtype.Text{String: "This is a bio", Valid: true},
		Company:     pgtype.Text{String: "company", Valid: true},
		JobRole:     "Developer",
		Description: pgtype.Text{String: "description", Valid: true},
	})
	assert.NoError(t, err)

	profileID := strconv.Itoa(int(profile.ID))

	payload := payloads.UpdateProfile{
		Bio:         pgtype.Text{String: "Updated bio", Valid: true},
		Company:     "Updated Company",
		JobRole:     "Senior Developer",
		Description: pgtype.Text{String: "Updated description", Valid: true},
		Skills:      []string{"Python", "Docker"},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPut, baseProfilePath+"/"+profileID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfully updated profile", response["status"])

	// Clean up
	truncateTables(t, "profiles", "users")
}

func TestGetProfileByID(t *testing.T) {
	router := setupTestRouter()

	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	profile, err := testQueries.CreateProfile(context.TODO(), sqlc.CreateProfileParams{
		UserID:      user.ID,
		Bio:         pgtype.Text{String: "This is a bio", Valid: true},
		Company:     pgtype.Text{String: "Example Company", Valid: true},
		JobRole:     "Developer",
		Description: pgtype.Text{String: "This is a description", Valid: true},
	})
	assert.NoError(t, err)

	profileID := strconv.Itoa(int(profile.ID))

	req, _ := http.NewRequest(http.MethodGet, baseProfilePath+"/"+profileID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrived profile", response["status"])

	// Clean up
	truncateTables(t, "profiles", "users")
}

func TestGetAllProfiles(t *testing.T) {
	router := setupTestRouter()

	// Create a user and a profile to ensure there's at least one profile in the database
	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	_, err = testQueries.CreateProfile(context.TODO(), sqlc.CreateProfileParams{
		UserID:      user.ID,
		Bio:         pgtype.Text{String: "This is a bio", Valid: true},
		Company:     pgtype.Text{String: "Example Company", Valid: true},
		JobRole:     "Developer",
		Description: pgtype.Text{String: "This is a description", Valid: true},
	})
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, baseProfilePath, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrieved all profile", response["status"])
	assert.Greater(t, int(response["size"].(float64)), 0)

	// Clean up
	truncateTables(t, "profiles", "users")
}

func TestDeleteProfileByID(t *testing.T) {
	router := setupTestRouter()

	hashedPassword, err := controllers.HashPassword("password123")
	assert.NoError(t, err)
	user, err := testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		FullName:     "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: hashedPassword,
		UserType:     "user",
	})
	assert.NoError(t, err)

	profile, err := testQueries.CreateProfile(context.TODO(), sqlc.CreateProfileParams{
		UserID:      user.ID,
		Bio:         pgtype.Text{String: "This is a bio", Valid: true},
		Company:     pgtype.Text{String: "Example Company", Valid: true},
		JobRole:     "Developer",
		Description: pgtype.Text{String: "This is a description", Valid: true},
	})
	assert.NoError(t, err)

	profileID := strconv.Itoa(int(profile.ID))

	req, _ := http.NewRequest(http.MethodDelete, baseProfilePath+"/"+profileID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfuly deleted", response["status"])

	// Clean up
	truncateTables(t, "profiles", "users")
}
