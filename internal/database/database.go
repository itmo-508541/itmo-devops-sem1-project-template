package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	*pgx.Conn
}

func New(ctx context.Context, dsn string) (*Database, error) {
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect: %w", err)
	}

	return &Database{db}, nil
}

func (d *Database) WithTransaction(
	parentCtx context.Context,
	do TransactionCallback,
) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	rawTx, err := d.Begin(ctx)
	if err != nil {
		return fmt.Errorf("d.Begin: %w", err)
	}

	tx := &transaction{rawTx}
	if err = do(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Println(fmt.Errorf("tx.Rollback: %w", err).Error())
		}

		return fmt.Errorf("do: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
