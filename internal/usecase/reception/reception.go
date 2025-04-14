package reception

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"pvz-service/internal/entity"
	"pvz-service/internal/repo"
	"pvz-service/internal/repo/postgres"
)

type UseCase struct {
	receptionRepo repo.ReceptionRepo
	trm           trm.Manager
}

func New(receptionRepo repo.ReceptionRepo, trm trm.Manager) *UseCase {
	return &UseCase{
		receptionRepo: receptionRepo,
		trm:           trm,
	}
}

func (uc *UseCase) Create(ctx context.Context, pvzId string) (*entity.Reception, error) {
	receptionNew, err := uc.receptionRepo.Create(ctx, pvzId)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateKey) {
			return nil, entity.ErrUnclosedReception
		}
		return nil, fmt.Errorf("ReceptionUseCase - Create - uc.receptionRepo.Create: %w", err)
	}

	return receptionNew, nil
}

func (uc *UseCase) GetByPvzId(ctx context.Context, pvzId string) ([]entity.Reception, error) {
	receptions, err := uc.receptionRepo.GetByPvzId(ctx, pvzId)
	if err != nil {
		return nil, fmt.Errorf("ReceptionUseCase - GetByPvzId - uc.receptionRepo.GetByPvzId: %w", err)
	}

	return receptions, nil
}

func (uc *UseCase) CloseLastOpenOnPvz(ctx context.Context, pvzId string) (*entity.Reception, error) {
	reception := entity.Reception{}

	err := uc.trm.Do(ctx, func(ctx context.Context) error {
		receptions, err := uc.receptionRepo.GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress)
		if err != nil {
			return fmt.Errorf("ReceptionUseCase - CloseLastOpenOnPvz - uc.receptionRepo.GetByPvzIdWithStatus: %w", err)
		}

		if len(receptions) == 0 {
			return entity.ErrNoOpenReception
		}

		reception = receptions[0]

		err = uc.receptionRepo.UpdateStatus(ctx, reception.ID, entity.Close)
		if err != nil {
			return fmt.Errorf("ReceptionUseCase - CloseLastOpenOnPvz - uc.receptionRepo.UpdateStatus: %w", err)
		}

		reception.Status = entity.Close

		return nil

	})

	if err != nil {
		return nil, err
	}

	return &reception, nil
}
