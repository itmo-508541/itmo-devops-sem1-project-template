package price

import (
	"context"
	"fmt"
	"project_sem/internal/database"
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	repository := &Repository{db: db}

	return repository
}

func (r *Repository) DeleteAll(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if err := r.db.WithTransaction(func(conn database.Connection) error {
		if _, err := conn.Exec(ctx, "DELETE FROM prices"); err != nil {
			return fmt.Errorf("conn.Exec: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("r.db.WithTransaction: %w", err)
	}

	return nil
}
