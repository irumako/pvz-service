package postgres

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"pvz-service/internal/entity"
	"pvz-service/pkg/postgres"
)

type ProductRepo struct {
	db     *postgres.Postgres
	getter *trmpgx.CtxGetter
}

func NewProductRepo(db *postgres.Postgres, getter *trmpgx.CtxGetter) *ProductRepo {
	return &ProductRepo{
		db:     db,
		getter: getter,
	}
}

func (r *ProductRepo) GetByReceptionId(ctx context.Context, receptionId string) ([]entity.Product, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Select("p.id", "p.datetime", "p.type").
		From("product p").
		Where("p.reception_id = ?", receptionId)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - GetByReceptionId - builder: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - GetByReceptionId - conn.Query: %w", err)
	}
	defer rows.Close()

	products := make([]entity.Product, 0)

	for rows.Next() {
		p := entity.Product{}
		p.ReceptionID = receptionId

		err = rows.Scan(&p.ID, &p.Datetime, &p.Type)
		if err != nil {
			return nil, fmt.Errorf("ProductRepo - GetByReceptionId - rows.Scan: %w", err)
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepo) Create(
	ctx context.Context,
	receptionId string,
	productType entity.ProductType,
) (*entity.Product, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Insert("product").
		Columns("reception_id", "type").
		Values(receptionId, productType).
		Suffix("RETURNING id, datetime")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - Create - builder: %w", err)
	}

	product := entity.Product{}
	product.ReceptionID = receptionId
	product.Type = productType

	err = conn.QueryRow(ctx, query, args...).Scan(&product.ID, &product.Datetime)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - Create - conn.QueryRow: %w", err)
	}

	return &product, nil
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db.Pool)
	builder := r.db.Builder.
		Delete("product").
		Where("id = ?", id)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("ProductRepo - Delete - builder: %w", err)
	}

	cmdTag, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ProductRepo - Delete - conn.Exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("ProductRepo - Delete - rowsAffected == 0: %s", id)
	}

	return nil
}
