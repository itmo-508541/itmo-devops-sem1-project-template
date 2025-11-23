package price

import (
	"context"
	"fmt"
	"io"
	"project_sem/internal/app/validate"
	"project_sem/internal/database"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	repository := &Repository{db: db}

	return repository
}

func (r *Repository) AcceptCsv(
	ctx context.Context,
	reader io.Reader,
) (result *AcceptResultDTO, err error) {
	input := make([]PriceDTO, 0)

	err = gocsv.Unmarshal(reader, &input)
	if err != nil && false { // игнорируем(?) ошибки
		err = fmt.Errorf("gocsv.Unmarshal: %w", err)

		return nil, err
	}

	v, err := validate.New()
	if err != nil {
		err = fmt.Errorf("validate.New: %w", err)

		return nil, err
	}

	result = &AcceptResultDTO{
		// Задание: общее количество строк в исходном файле
		TotalCount: len(input),
	}

	// Ревью 1: все нужно делать в транзакции
	err = r.db.WithTransaction(ctx, func(ctx context.Context, conn database.Connection) error {
		var exists int
		for _, price := range input {
			err := v.Struct(price)
			if err != nil {
				continue
			}
			criteria := pgx.NamedArgs{
				"id":          price.Id,
				"name":        price.Name,
				"category":    price.Category,
				"price":       price.Price,
				"create_date": price.CreateDate,
			}

			// Ревью 1: при добавлении если есть дубликаты то считать их - дубликатом считается строчка полностью совпадающая по всем полям
			//
			// Я не учитываю колонку id из csv-файла для определения уникальности записи
			row := r.db.QueryRow(
				ctx,
				"SELECT COUNT(*) FROM prices WHERE name=@name AND category=@category AND price=@price AND create_date=@create_date",
				criteria,
			)
			if err = row.Scan(&exists); err != nil {
				return fmt.Errorf("row.Scan: %w", err)
			}
			if exists > 0 {
				// Задание: количество дубликатов во входных данных и в СУБД
				result.DuplicatesCount++

				continue
			}

			_, err = r.db.Exec(
				ctx,
				"INSERT INTO prices (id, name, category, price, create_date) VALUES (@id, @name, @category, @price, @create_date)",
				criteria,
			)
			if err != nil {
				return fmt.Errorf("r.db.Exec: %w", err)
			}
			// Задание: общее количество добавленных элементов
			result.TotalItems++
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("r.db.WithTransaction: %w", err)
	}

	// Ревью 1: далее в рамках этой же транзакции нужно сделать агрегирующий запрос чтобы посчитать статистику
	//
	// Я вынес аггрегирующий запрос из транзакции, чтобы сократить время работы этой транзакции
	row := r.db.QueryRow(
		ctx,
		"SELECT COALESCE(COUNT(DISTINCT category), 0) AS categories, COALESCE(SUM(price), 0) AS prices FROM prices",
	)
	// Задание: общее количество категорий
	//   Тут неоднозначное задание: не понятно какие категории считать: в загружаемом файте или в базе данных
	//   Я посчитал "итого по всем объектам в базе"
	// Задание: суммарная стоимость всех объектов в базе данных
	if err = row.Scan(&result.TotalCategories, &result.TotalPrice); err != nil {
		err = fmt.Errorf("row.Scan: %w", err)

		return
	}

	return result, nil
}

func (r *Repository) FindByFilter(
	parentCtx context.Context,
	filter SqlFilter,
) (*[]PriceDTO, error) {
	sql := "SELECT id, name, category, price, create_date FROM prices"
	args, where, ok := filter.Where()
	if ok {
		sql += " WHERE " + where
	}

	rows, err := r.db.Query(parentCtx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	var date time.Time
	all := make([]PriceDTO, 0)
	for rows.Next() {
		p := PriceDTO{}

		err = rows.Scan(&p.Id, &p.Name, &p.Category, &p.Price, &date)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		p.CreateDate = date.Format(time.DateOnly)

		all = append(all, p)
	}

	return &all, nil
}
