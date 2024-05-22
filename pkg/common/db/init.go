package db

import (
	"context"
	"errors"

	"client_task/pkg/common/config"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

// Config holds the database connection configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// Connect creates a new database connection
func Connect(ctx context.Context) (*pgx.Conn, error) {
	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		Database: viper.GetString("POSTGRES_DB"),
	}

	return connectToDatabase(ctx, cfg)
}

// connectToDatabase establishes the database connection (private function)
func connectToDatabase(ctx context.Context, cfg Config) (*pgx.Conn, error) {

	connString := "postgres://" + cfg.User + ":" + cfg.Password + "@" +
		cfg.Host + ":" + cfg.Port + "/" + cfg.Database

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, errors.New("failed to ping database: " + err.Error())
	}

	return conn, nil
}
