package postgres

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"pvz-service/internal/entity"
	"pvz-service/pkg/postgres"
)

type ReceptionRepo struct {
	db     *postgres.Postgres
	getter *trmpgx.CtxGetter
}

func NewReceptionRepo(db *postgres.Postgres, getter *trmpgx.CtxGetter) *ReceptionRepo {
	return &ReceptionRepo{
		db:     db,
		getter: getter,
	}
}

func (r *ReceptionRepo) GetById(ctx context.Context, id string) (*entity.Reception, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Select("r.datetime", "r.pvz_id", "r.status").
		From("reception r").
		Where("r.id = ?", id)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - GetById - builder: %w", err)
	}

	reception := entity.Reception{}

	err = conn.QueryRow(ctx, query, args...).Scan(&reception.Datetime, &reception.PvzID, &reception.Status)
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - GetById - conn.QueryRow: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepo) GetByPvzId(ctx context.Context, pvzId string) ([]entity.Reception, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Select("r.id", "r.datetime", "r.status").
		From("reception r").
		Where("r.pvz_id = ?", pvzId)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - GetByPvzId - builder: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - GetByPvzId - conn.Query: %w", err)
	}
	defer rows.Close()

	receptions := make([]entity.Reception, 0)

	for rows.Next() {
		r := entity.Reception{}
		r.PvzID = pvzId

		err = rows.Scan(&r.ID, &r.Datetime, &r.Status)
		if err != nil {
			return nil, fmt.Errorf("ReceptionRepo - GetByPvzId - rows.Scan: %w", err)
		}

		receptions = append(receptions, r)
	}

	return receptions, nil
}

func (r *ReceptionRepo) GetByPvzIdWithStatus(
	ctx context.Context,
	pvzId string,
	status entity.ReceptionStatus,
) ([]entity.Reception, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Select("r.id", "r.datetime").
		From("reception r").
		Where("r.pvz_id = ?", pvzId).
		Where("r.status = ?", status)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - GetByPvzIdWithStatus - builder: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - GetByPvzIdWithStatus - conn.Query: %w", err)
	}
	defer rows.Close()

	receptions := make([]entity.Reception, 0)

	for rows.Next() {
		r := entity.Reception{}
		r.PvzID = pvzId
		r.Status = status

		err = rows.Scan(&r.ID, &r.Datetime)
		if err != nil {
			return nil, fmt.Errorf("ReceptionRepo - GetByPvzIdWithStatus - rows.Scan: %w", err)
		}

		receptions = append(receptions, r)
	}

	return receptions, nil
}

func (r *ReceptionRepo) Create(ctx context.Context, pvzId string) (*entity.Reception, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Insert("reception").
		Columns("pvz_id").
		Values(pvzId).
		Suffix("RETURNING id, datetime, status")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - Create - builder: %w", err)
	}

	reception := entity.Reception{}
	reception.PvzID = pvzId

	err = conn.QueryRow(ctx, query, args...).Scan(&reception.ID, &reception.Datetime, &reception.Status)
	if err != nil {
		return nil, fmt.Errorf("ReceptionRepo - Create - conn.QueryRow: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepo) UpdateStatus(ctx context.Context, id string, status entity.ReceptionStatus) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Update("reception").
		Set("status", status).
		Where("id = ?", id).
		Suffix("RETURNING datetime, pvz_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("ReceptionRepo - UpdateStatus - builder: %w", err)
	}

	cmdTag, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ReceptionRepo - UpdateStatus - conn.Exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("ReceptionRepo - UpdateStatus - rowsAffected == 0: %s", id)
	}

	return nil
}
