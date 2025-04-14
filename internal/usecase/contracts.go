package usecase

import (
	"context"
	"pvz-service/internal/entity"
	"time"
)

type (
	Pvz interface {
		Create(ctx context.Context, pvz entity.Pvz) (*entity.Pvz, error)
		ListByReceptionDate(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entity.Pvz, error)
	}
	Reception interface {
		Create(ctx context.Context, pvzId string) (*entity.Reception, error)
		GetByPvzId(ctx context.Context, pvzId string) ([]entity.Reception, error)
		CloseLastOpenOnPvz(ctx context.Context, pvzId string) (*entity.Reception, error)
	}
	Product interface {
		GetByReceptionId(ctx context.Context, receptionId string) ([]entity.Product, error)
		Create(ctx context.Context, pvzId string, productType entity.ProductType) (*entity.Product, error)
		Delete(ctx context.Context, pvzId string) error
	}
	User interface {
	}
)
