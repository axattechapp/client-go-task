package main

import (
	"client_task/pkg/common/config"
	"client_task/pkg/common/db"
	sqlc "client_task/pkg/common/db/sqlc"
	"client_task/pkg/router"
	"context"
	"fmt"
)

func main() {
	fmt.Println("program started...")
	config.LoadConfig()
	ctx := context.Background()

	// Connect to the database
	conn, err := db.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	fmt.Println("Successfully connected to database!")

	defer conn.Close(ctx) // Close the connection when the program exits

	// Use the connection object here for database operations
	db := sqlc.New(conn)

	// Setup the router
	server := router.SetupRouter(db, ctx)

	// Start the server
	server.Run(":8080")
}
