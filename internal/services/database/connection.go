package database

import "github.com/jackc/pgx/v5"

type Connection struct {
	*pgx.Conn
}
