package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type transaction struct {
	pgx.Tx
}

func (tx *transaction) WithTransaction(ctx context.Context, do TransactionCallback) error {
	return do(ctx, tx)
}
