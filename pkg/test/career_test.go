package test

import (
	"bytes"
	sqlc "client_task/pkg/common/db/sqlc"
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

var baseCareerPath string = "/api/careers"

func TestCreateCareer(t *testing.T) {
	router := setupTestRouter()

	payload := payloads.CreateCareer{
		UserID:      1,
		Title:       "Software Engineer",
		Company:     "Tech Corp",
		Description: "Developing software solutions",
		SkillID:     1,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseCareerPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfully created career", response["status"])

	// Clean up
	truncateTables(t, "careers")
}

func TestUpdateCareer(t *testing.T) {
	router := setupTestRouter()

	// Create a career first
	career, err := testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      1,
		Title:       "Software Engineer",
		Company:     pgtype.Text{String: "Tech Corp", Valid: true},
		Description: pgtype.Text{String: "Developing software solutions", Valid: true},
		SkillID:     1,
	})
	assert.NoError(t, err)

	payload := payloads.UpdateCareer{
		Title:       "Senior Software Engineer",
		Company:     "Tech Corp",
		Description: "Leading software projects",
		SkillID:     2,
	}
	body, _ := json.Marshal(payload)

	careerID := strconv.Itoa(int(career.ID))
	req, _ := http.NewRequest(http.MethodPut, baseCareerPath+"/"+careerID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "successfully updated career", response["status"])

	// Clean up
	truncateTables(t, "careers")
}

func TestGetCareerByID(t *testing.T) {
	router := setupTestRouter()

	// Create a career first
	career, err := testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      1,
		Title:       "Software Engineer",
		Company:     pgtype.Text{String: "Tech Corp", Valid: true},
		Description: pgtype.Text{String: "Developing software solutions", Valid: true},
		SkillID:     1,
	})
	assert.NoError(t, err)

	careerID := strconv.Itoa(int(career.ID))
	req, _ := http.NewRequest(http.MethodGet, baseCareerPath+"/"+careerID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrieved career", response["status"])

	// Clean up
	truncateTables(t, "careers")
}

func TestGetAllCareers(t *testing.T) {
	router := setupTestRouter()

	// Create some careers first
	_, err := testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      1,
		Title:       "Software Engineer",
		Company:     pgtype.Text{String: "Tech Corp", Valid: true},
		Description: pgtype.Text{String: "Developing software solutions", Valid: true},
		SkillID:     1,
	})
	assert.NoError(t, err)
	_, err = testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      2,
		Title:       "Data Scientist",
		Company:     pgtype.Text{String: "Data Corp", Valid: true},
		Description: pgtype.Text{String: "Analyzing data", Valid: true},
		SkillID:     2,
	})
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, baseCareerPath, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrieved all careers", response["status"])
	assert.Greater(t, int(response["size"].(float64)), 0)

	// Clean up
	truncateTables(t, "careers")
}

func TestDeleteCareerByID(t *testing.T) {
	router := setupTestRouter()

	// Create a career first
	career, err := testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      1,
		Title:       "Software Engineer",
		Company:     pgtype.Text{String: "Tech Corp", Valid: true},
		Description: pgtype.Text{String: "Developing software solutions", Valid: true},
		SkillID:     1,
	})
	assert.NoError(t, err)

	careerID := strconv.Itoa(int(career.ID))
	req, _ := http.NewRequest(http.MethodDelete, baseCareerPath+"/"+careerID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully deleted career", response["status"])

	// Clean up
	truncateTables(t, "careers")
}

func TestGetCareersByUserID(t *testing.T) {
	router := setupTestRouter()

	// Create some careers first
	_, err := testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      1,
		Title:       "Software Engineer",
		Company:     pgtype.Text{String: "Tech Corp", Valid: true},
		Description: pgtype.Text{String: "Developing software solutions", Valid: true},
		SkillID:     1,
	})
	assert.NoError(t, err)
	_, err = testQueries.CreateCareer(context.TODO(), sqlc.CreateCareerParams{
		UserID:      1,
		Title:       "Data Scientist",
		Company:     pgtype.Text{String: "Data Corp", Valid: true},
		Description: pgtype.Text{String: "Analyzing data", Valid: true},
		SkillID:     2,
	})
	assert.NoError(t, err)

	userID := "1"
	req, _ := http.NewRequest(http.MethodGet, baseCareerPath+"/user/"+userID, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrieved careers", response["status"])
	assert.Greater(t, int(response["size"].(float64)), 0)

	// Clean up
	truncateTables(t, "careers")
}
