package usecase

import (
	"context"
	"github.com/stretchr/testify/require"
	"pvz-service/internal/entity"
	"pvz-service/internal/repo/postgres"
	"pvz-service/internal/usecase/product"
	"testing"
	"time"
)

func productUseCase(t *testing.T) (*product.UseCase, *MockProductRepo, *MockReceptionRepo) {
	t.Helper()

	productRepo := NewMockProductRepo(t)
	receptionRepo := NewMockReceptionRepo(t)
	fakeManager := &FakeManager{}

	uc := product.New(productRepo, receptionRepo, fakeManager)
	return uc, productRepo, receptionRepo
}

func TestProductUseCase_GetByReceptionId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pvzId := "pvz-uuid4"
	fixedTime, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	productType := entity.Electronic
	openReception := entity.Reception{
		ID:       "rec-001",
		PvzID:    pvzId,
		Datetime: fixedTime,
		Status:   entity.InProgress,
	}
	productsExpected := []entity.Product{
		{
			ID:          "prod-001",
			ReceptionID: openReception.ID,
			Datetime:    fixedTime,
			Type:        productType,
		},
		{
			ID:          "prod-002",
			ReceptionID: openReception.ID,
			Datetime:    fixedTime,
			Type:        productType,
		},
	}

	tests := []struct {
		name string
		mock func(productRepo *MockProductRepo, recRepo *MockReceptionRepo)
		res  []entity.Product
		err  error
	}{
		{
			name: "Успешное получение товаров",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {

				productRepo.EXPECT().GetByReceptionId(ctx, openReception.ID).Return(productsExpected, nil)
			},
			res: productsExpected,
			err: nil,
		},
		{
			name: "Ошибка в репо",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				productRepo.EXPECT().GetByReceptionId(ctx, openReception.ID).Return(nil, entity.ErrInternalServErr)
			},
			res: nil,
			err: entity.ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			uc, productRepo, receptionRepo := productUseCase(t)

			localTc.mock(productRepo, receptionRepo)

			res, err := uc.GetByReceptionId(ctx, openReception.ID)

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestProductUseCase_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pvzId := "pvz-uuid4"
	fixedTime, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	productType := entity.Electronic
	openReception := entity.Reception{
		ID:       "rec-001",
		PvzID:    pvzId,
		Datetime: fixedTime,
		Status:   entity.InProgress,
	}

	tests := []struct {
		name string
		mock func(productRepo *MockProductRepo, recRepo *MockReceptionRepo)
		res  *entity.Product
		err  error
	}{
		{
			name: "Успешное создание товара",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				recRepo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{openReception}, nil)
				prod := &entity.Product{
					ID:          "prod-001",
					ReceptionID: openReception.ID,
					Datetime:    fixedTime,
					Type:        productType,
				}
				productRepo.EXPECT().Create(ctx, openReception.ID, productType).Return(prod, nil)
			},
			res: &entity.Product{
				ID:          "prod-001",
				ReceptionID: openReception.ID,
				Datetime:    fixedTime,
				Type:        productType,
			},
			err: nil,
		},
		{
			name: "Нет открытой приёмки",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				recRepo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{}, nil)
			},
			res: nil,
			err: entity.ErrNoOpenReception,
		},
		{
			name: "Ошибка в репо при создании",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				recRepo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{openReception}, nil)
				productRepo.EXPECT().
					Create(ctx, openReception.ID, productType).
					Return(nil, entity.ErrInternalServErr)
			},
			res: nil,
			err: entity.ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			uc, productRepo, receptionRepo := productUseCase(t)

			localTc.mock(productRepo, receptionRepo)

			res, err := uc.Create(ctx, pvzId, productType)

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestProductUseCase_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pvzId := "pvz-uuid4"
	fixedTime, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	productType := entity.Electronic
	openReception := entity.Reception{
		ID:       "rec-001",
		PvzID:    pvzId,
		Datetime: fixedTime,
		Status:   entity.InProgress,
	}
	lastProduct := &entity.Product{
		ID:          "prod-001",
		ReceptionID: openReception.ID,
		Datetime:    fixedTime,
		Type:        productType,
	}

	tests := []struct {
		name string
		mock func(productRepo *MockProductRepo, recRepo *MockReceptionRepo)
		err  error
	}{
		{
			name: "Успешное удаление товара",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				recRepo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{openReception}, nil)

				productRepo.EXPECT().
					GetLastInReception(ctx, openReception.ID).
					Return(lastProduct, nil)

				productRepo.EXPECT().
					Delete(ctx, lastProduct.ID).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "Нет открытой приёмки",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				recRepo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{}, nil)
			},
			err: entity.ErrNoOpenReception,
		},
		{
			name: "Отсутствуют товары",
			mock: func(productRepo *MockProductRepo, recRepo *MockReceptionRepo) {
				recRepo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{openReception}, nil)
				productRepo.EXPECT().
					GetLastInReception(ctx, openReception.ID).
					Return(nil, postgres.ErrRecordNotFound)
			},
			err: entity.ErrNoProductsToDelete,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			uc, productRepo, receptionRepo := productUseCase(t)

			localTc.mock(productRepo, receptionRepo)

			err = uc.Delete(ctx, pvzId)

			require.ErrorIs(t, err, localTc.err)
		})
	}
}
