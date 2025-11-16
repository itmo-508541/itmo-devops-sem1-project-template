package price

import (
	"context"
	"fmt"
	"io"
	"project_sem/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

type Repository struct {
	db        *database.Database
	validator *validator.Validate
}

func NewRepository(db *database.Database, v *validator.Validate) *Repository {
	repository := &Repository{db: db, validator: v}

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

func (r *Repository) AcceptCsv(
	parentCtx context.Context,
	reader io.Reader,
) (uid uuid.UUID, totalCount int, err error) {
	input := make([]PriceRecordDTO, 0)
	output := make([]PriceRecordDTO, 0)

	err = gocsv.Unmarshal(reader, &input)
	if err != nil && false { // игнорируем(?) ошибки
		err = fmt.Errorf("gocsv.Unmarshal: %w", err)

		return
	}

	uid = uuid.New()
	for _, price := range input {
		err := r.validator.Struct(price)
		if err != nil {
			continue
		}
		price.GroupUUID = uid
		price.UUID = uuid.New()
		output = append(output, price)
	}

	return uid, len(input), r.insertAll(parentCtx, &output)
}

func (r *Repository) insertAll(parentCtx context.Context, prices *[]PriceRecordDTO) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	for _, price := range *prices {
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
			return fmt.Errorf("db.Exec: %w", err)
		}
	}

	return nil
}
