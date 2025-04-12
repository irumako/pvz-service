package postgres

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"pvz-service/internal/entity"
	"pvz-service/pkg/postgres"
	"time"
)

type PvzRepo struct {
	db     *postgres.Postgres
	getter *trmpgx.CtxGetter
}

func NewPvzRepo(db *postgres.Postgres, getter *trmpgx.CtxGetter) *PvzRepo {
	return &PvzRepo{
		db:     db,
		getter: getter,
	}
}

func (r *PvzRepo) GetById(ctx context.Context, id string) (*entity.Pvz, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Select("pvz.id", "pvz.registration_date", "pvz.city").
		From("pvz").
		Where("pvz.id = ?", id)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PvzRepo - GetById - builder: %w", err)
	}

	pvz := entity.Pvz{}

	err = conn.QueryRow(ctx, query, args...).Scan(&pvz.ID, &pvz.RegistrationDate)
	if err != nil {
		return nil, fmt.Errorf("PvzRepo - GetById - conn.QueryRow: %w", err)
	}

	return &pvz, nil
}

func (r *PvzRepo) GetByReceptionDate(
	ctx context.Context,
	startDate, endDate *time.Time,
	limit, offset int,
) ([]entity.Pvz, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Select("pvz.id", "pvz.registration_date", "pvz.city").
		From("pvz").
		Join("reception r ON pvz.id = r.pvz_id")

	if startDate != nil {
		builder = builder.Where("r.reception_datetime >= ?", *startDate)
	}
	if endDate != nil {
		builder = builder.Where("r.reception_datetime <= ?", *endDate)
	}

	builder = builder.
		OrderBy("pvz.registration_date ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PvzRepo - GetByReceptionDate - builder: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PvzRepo - GetByReceptionDate - conn.Query: %w", err)
	}
	defer rows.Close()

	pvzList := make([]entity.Pvz, 0)

	for rows.Next() {
		p := entity.Pvz{}

		err = rows.Scan(&p.ID, &p.City, &p.RegistrationDate)
		if err != nil {
			return nil, fmt.Errorf("PvzRepo - GetByReceptionDate - rows.Scan: %w", err)
		}

		pvzList = append(pvzList, p)
	}

	return pvzList, nil
}

func (r *PvzRepo) Create(ctx context.Context, pvz entity.Pvz) (*entity.Pvz, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Insert("pvz").
		Columns("id, registration_date, city").
		Values(pvz.ID, pvz.RegistrationDate, pvz.City)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PvzRepo - Create - builder: %w", err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PvzRepo - Create - conn.Exec: %w", err)
	}

	return &pvz, nil
}
