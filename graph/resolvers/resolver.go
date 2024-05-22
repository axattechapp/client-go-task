package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import db_sqlc "client_task/pkg/common/db/sqlc"

//go:generate go get github.com/99designs/gqlgen/cmd@v0.14.0
//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	DB *db_sqlc.Queries
}
