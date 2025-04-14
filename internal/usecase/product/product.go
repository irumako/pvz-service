package product

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
	productRepo   repo.ProductRepo
	receptionRepo repo.ReceptionRepo
	trm           trm.Manager
}

func New(productRepo repo.ProductRepo, receptionRepo repo.ReceptionRepo, trm trm.Manager) *UseCase {
	return &UseCase{
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
		trm:           trm,
	}
}

func (uc *UseCase) Create(ctx context.Context, pvzId string, productType entity.ProductType) (*entity.Product, error) {

	product := &entity.Product{}

	err := uc.trm.Do(ctx, func(ctx context.Context) error {
		receptions, err := uc.receptionRepo.GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress)
		if err != nil {
			return fmt.Errorf("ProductUseCase - Create - uc.receptionRepo.GetByPvzIdWithStatus: %w", err)
		}

		if len(receptions) == 0 {
			return entity.ErrNoOpenReception
		}

		reception := receptions[0]

		product, err = uc.productRepo.Create(ctx, reception.ID, productType)
		if err != nil {
			return fmt.Errorf("ProductUseCase - Create - uc.productRepo.Create: %w", err)
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *UseCase) Delete(ctx context.Context, pvzId string) error {
	err := uc.trm.Do(ctx, func(ctx context.Context) error {
		receptions, err := uc.receptionRepo.GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress)
		if err != nil {
			return fmt.Errorf("ProductUseCase - Delete - uc.receptionRepo.GetByPvzIdWithStatus: %w", err)
		}

		if len(receptions) == 0 {
			return entity.ErrNoOpenReception
		}

		reception := receptions[0]

		lastProduct, err := uc.productRepo.GetLastInReception(ctx, reception.ID)
		if err != nil {
			if errors.Is(err, postgres.ErrRecordNotFound) {
				return entity.ErrNoProductsToDelete
			}
			return fmt.Errorf("ProductUseCase - Delete - uc.productRepo.GetLastInReception: %w", err)
		}

		if err := uc.productRepo.Delete(ctx, lastProduct.ID); err != nil {
			return fmt.Errorf("ProductUseCase - Delete - uc.productRepo.Delete: %w", err)
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}
