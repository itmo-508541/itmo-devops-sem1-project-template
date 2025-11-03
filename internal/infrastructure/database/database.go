package database

import (
	"database/sql"
	"fmt"
	"os"
	"project_sem/internal/env"
)

type Database struct {
	Connection *sql.DB
}

func New() (*Database, error) {
	connection, err := sql.Open("postgres", DataSourceName())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return &Database{Connection: connection}, err
}

// @todo поменять =) на то, что в вебинаре!
func DataSourceName() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv(env.DatabaseUser),
		os.Getenv(env.DatabasePassword),
		os.Getenv(env.DatabaseHost),
		os.Getenv(env.DatabasePort),
		os.Getenv(env.DatabaseName),
	)
}
