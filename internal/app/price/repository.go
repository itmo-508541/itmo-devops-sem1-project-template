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

	if _, err := r.db.Exec(ctx, "DELETE FROM prices"); err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}

func (r *Repository) InsertAll(parentCtx context.Context, prices *[]PriceRecordDTO) error {
	for _, price := range *prices {
		err := r.Insert(parentCtx, &price)
		if err != nil {
			return fmt.Errorf("Insert: %w", err)
		}
	}

	return nil
}

func (r *Repository) Insert(parentCtx context.Context, price *PriceRecordDTO) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if _, err := r.db.Exec(
		ctx,
		"INSERT INTO prices (group_uuid, uuid, id, name, category, price, create_date) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		price.GroupUUID.String(),
		price.UUID.String(),
		price.Id, price.Name,
		price.Category,
		price.Price,
		price.CreateDate,
	); err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}
