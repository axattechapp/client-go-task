package test

import (
	"bytes"
	"client_task/pkg/payloads"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseSkillPath string = "/api/skills"

func TestCreateSkill(t *testing.T) {
	router := setupTestRouter()

	payload := payloads.CreateSkill{
		Name: "Python",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseSkillPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully created skill", response["status"])

	// Clean up
	truncateTables(t, "skills")
}

func TestGetAllSkills(t *testing.T) {
	router := setupTestRouter()

	// Create some skills to ensure there are skills in the database
	_, err := testQueries.CreateSkill(context.TODO(), "JavaScript")
	assert.NoError(t, err)
	_, err = testQueries.CreateSkill(context.TODO(), "Java")
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, baseSkillPath, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully retrieved all skills", response["status"])
	assert.Greater(t, int(response["size"].(float64)), 0)

	// Clean up
	truncateTables(t, "skills")
}
