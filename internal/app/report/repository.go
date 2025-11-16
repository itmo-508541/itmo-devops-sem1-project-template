package report

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"project_sem/internal/database"
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	repository := &Repository{db: db}

	return repository
}

func (r *Repository) Renew(
	parentCtx context.Context,
	UUID uuid.UUID,
) (duplicatesCount int, totalItems int, totalCategories int, totalPrice float32, err error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if _, err = r.db.Exec(ctx, "DELETE FROM reports"); err != nil {
		err = fmt.Errorf("db.Exec(DELETE reports): %w", err)
		return
	}

	sql := `
INSERT INTO reports
	(uuid, id, name, category, price, create_date)
SELECT
	MIN(uuid), MIN(id), name, category, price, create_date
FROM
	prices
WHERE
	group_uuid=$1
GROUP BY 
	name, category, price, create_date`

	if _, err = r.db.Exec(ctx, sql, UUID.String()); err != nil {
		err = fmt.Errorf("db.Exec(INSERT INTO reports): %w", err)
		return
	}

	var inserted, items, categories int
	var prices float32
	if err = r.db.QueryRow(ctx, "SELECT COALESCE(COUNT(*), 0) FROM prices WHERE group_uuid=$1", UUID.String()).Scan(&inserted); err != nil {
		err = fmt.Errorf("row.Scan: %w", err)
		return
	}

	if err = r.db.QueryRow(ctx, "SELECT COALESCE(COUNT(*), 0) AS items, COALESCE(COUNT(DISTINCT category), 0) AS categories, COALESCE(SUM(price), 0) AS prices FROM reports").Scan(&items, &categories, &prices); err != nil {
		err = fmt.Errorf("row.Scan: %w", err)
		return
	}

	return inserted - items, items, categories, prices, nil
}

func (r *Repository) All(parentCtx context.Context, filter RequestFilter) (*[]ReportDTO, error) {
	sql := "SELECT id, name, category, price, create_date FROM reports"
	args, where, ok := filter.Where()
	if ok {
		sql += fmt.Sprintf(" WHERE %s", where)
	}
	// fmt.Println(sql)
	// fmt.Println(args)

	rows, err := r.db.Query(parentCtx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	var date time.Time
	all := make([]ReportDTO, 0)
	for rows.Next() {
		p := ReportDTO{}

		err = rows.Scan(&p.Id, &p.Name, &p.Category, &p.Price, &date)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		p.CreateDate = date.Format(time.DateOnly)

		all = append(all, p)
	}

	return &all, nil
}
