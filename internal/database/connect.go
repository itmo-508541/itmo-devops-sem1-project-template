package database

import (
	"context"

	"github.com/jackc/pgx/v5"

	"fmt"
)

type Connection struct {
	*pgx.Conn
}

func New(dsn string) (*Connection, error) {
	connection, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect: %w", err)
	}

	return &Connection{connection}, err
}
