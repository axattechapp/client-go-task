// test/router_test.go
package test

import (
	"client_task/pkg/common/db"
	sqlc "client_task/pkg/common/db/sqlc"
	"client_task/pkg/router"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	ctx := context.Background()

	// Use the same DB connection setup as in your application
	conn, err := db.Connect(ctx)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	db := sqlc.New(conn)

	return router.SetupRouter(db, ctx)
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/healthcheck", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "The contact API is working fine")
}
