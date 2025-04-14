package pvz

import (
	"context"
	"fmt"
	"pvz-service/internal/entity"
	"pvz-service/internal/repo"
	"time"
)

type UseCase struct {
	pvzRepo repo.PvzRepo
}

func New(pvzRepo repo.PvzRepo) *UseCase {
	return &UseCase{
		pvzRepo: pvzRepo,
	}
}

func (uc *UseCase) Create(ctx context.Context, pvz entity.Pvz) (*entity.Pvz, error) {
	pvzNew, err := uc.pvzRepo.Create(ctx, pvz)
	if err != nil {
		return nil, fmt.Errorf("PvzUseCase - Create - uc.pvzRepo.Create: %w", err)
	}

	return pvzNew, nil
}

func (uc *UseCase) ListByReceptionDate(
	ctx context.Context,
	startDate, endDate *time.Time,
	page, limit int,
) ([]entity.Pvz, error) {
	pvzList, err := uc.pvzRepo.GetByReceptionDate(ctx, startDate, endDate, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("PvzUseCase - ListByReceptionDate - uc.pvzRepo.GetByReceptionDate: %w", err)
	}

	return pvzList, nil
}
