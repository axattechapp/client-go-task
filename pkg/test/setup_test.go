package test

import (
	sqlc "client_task/pkg/common/db/sqlc"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"client_task/pkg/common/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

var testQueries *sqlc.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.TODO()

	conn, err := db.Connect(ctx)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = sqlc.New(conn)

	code := m.Run()
	os.Exit(code)
}

func truncateTables(t *testing.T, tables ...string) {
	for _, table := range tables {
		_, err := testDB.ExecContext(context.Background(), "TRUNCATE TABLE "+table+" RESTART IDENTITY CASCADE")
		assert.NoError(t, err)
	}
}
