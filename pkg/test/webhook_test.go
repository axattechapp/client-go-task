package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleWebhook_JobApplication(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]interface{}{
		"eventType": "job_application",
		"jobId":     "123",
		"data": map[string]interface{}{
			"jobId":     "123",
			"applicant": "developer",
			"skills":    []interface{}{"Java", "Python"},
		},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestHandleWebhook_DailyJobs(t *testing.T) {
	return
}

func TestHandleWebhook_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	// Invalid JSON payload
	body := bytes.NewBuffer([]byte(`{"invalid": "json`))
	req, _ := http.NewRequest(http.MethodPost, "/webhook", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
