package repo

import (
	"context"
	"pvz-service/internal/entity"
	"time"
)

type (
	PvzRepo interface {
		GetById(ctx context.Context, id string) (*entity.Pvz, error)
		GetByReceptionDate(ctx context.Context, startDate, endDate *time.Time, limit, offset int) ([]entity.Pvz, error)
		Create(ctx context.Context, pvz entity.Pvz) (*entity.Pvz, error)
	}
	ReceptionRepo interface {
		GetById(ctx context.Context, id string) (*entity.Reception, error)
		GetByPvzId(ctx context.Context, pvzId string) ([]entity.Reception, error)
		GetByPvzIdWithStatus(ctx context.Context, pvzId string, status entity.ReceptionStatus) ([]entity.Reception, error)
		Create(ctx context.Context, pvzId string) (*entity.Reception, error)
		UpdateStatus(ctx context.Context, id string, status entity.ReceptionStatus) error
	}
	ProductRepo interface {
		GetByReceptionId(ctx context.Context, receptionId string) ([]entity.Product, error)
		GetLastInReception(ctx context.Context, receptionId string) (*entity.Product, error)
		Create(ctx context.Context, receptionId string, productType entity.ProductType) (*entity.Product, error)
		Delete(ctx context.Context, id string) error
	}
	UserRepo interface {
	}
)
