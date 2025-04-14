package reception

import (
	"context"
	"github.com/stretchr/testify/require"
	mocks "pvz-service/internal/repo"
	"pvz-service/internal/repo/postgres"
	"pvz-service/internal/usecase"
	"testing"
	"time"

	"pvz-service/internal/entity"
)

func receptionUseCase(t *testing.T) (*UseCase, *mocks.MockReceptionRepo) {
	t.Helper()

	repo := mocks.NewMockReceptionRepo(t)
	fakeManager := &usecase.FakeManager{}

	uc := New(repo, fakeManager)
	return uc, repo
}

func TestReceptionUseCase_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pvzId := "pvz-uuid4"
	fixedTime, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	tests := []struct {
		name string
		mock func(repo *mocks.MockReceptionRepo)
		res  *entity.Reception
		err  error
	}{
		{
			name: "Успешное создание приемки",
			mock: func(repo *mocks.MockReceptionRepo) {
				receptionFake := &entity.Reception{
					ID:       "uuid4",
					PvzID:    pvzId,
					Datetime: fixedTime,
					Status:   entity.InProgress,
				}
				repo.EXPECT().Create(ctx, pvzId).Return(receptionFake, nil)
			},
			res: &entity.Reception{
				ID:       "uuid4",
				PvzID:    pvzId,
				Datetime: fixedTime,
				Status:   entity.InProgress,
			},
			err: nil,
		},
		{
			name: "Невозможно создать приемку (уже есть открытая)",
			mock: func(repo *mocks.MockReceptionRepo) {
				repo.EXPECT().Create(ctx, pvzId).Return(nil, postgres.ErrDuplicateKey)
			},
			res: nil,
			err: entity.ErrUnclosedReception,
		},
		{
			name: "Ошибка в репо при создании приемки",
			mock: func(repo *mocks.MockReceptionRepo) {
				repo.EXPECT().Create(ctx, pvzId).Return(nil, entity.ErrInternalServErr)
			},
			res: nil,
			err: entity.ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			uc, repo := receptionUseCase(t)

			localTc.mock(repo)

			res, err := uc.Create(ctx, pvzId)

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestReceptionUseCase_GetByPvzId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pvzId := "pvz-uuid-456"

	fixedTime, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	tests := []struct {
		name string
		mock func(repo *mocks.MockReceptionRepo)
		res  []entity.Reception
		err  error
	}{
		{
			name: "Успешное получение приемок",
			mock: func(repo *mocks.MockReceptionRepo) {
				receptionsExpected := []entity.Reception{
					{ID: "rec1", PvzID: pvzId, Datetime: fixedTime, Status: entity.InProgress},
					{ID: "rec2", PvzID: pvzId, Datetime: fixedTime, Status: entity.Close},
				}
				repo.EXPECT().GetByPvzId(ctx, pvzId).Return(receptionsExpected, nil)
			},
			res: []entity.Reception{
				{ID: "rec1", PvzID: pvzId, Datetime: fixedTime, Status: entity.InProgress},
				{ID: "rec2", PvzID: pvzId, Datetime: fixedTime, Status: entity.Close},
			},
			err: nil,
		},
		{
			name: "Ошибка в репо",
			mock: func(repo *mocks.MockReceptionRepo) {
				repo.EXPECT().GetByPvzId(ctx, pvzId).Return(nil, entity.ErrInternalServErr)
			},
			res: nil,
			err: entity.ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			uc, repo := receptionUseCase(t)

			localTc.mock(repo)

			res, err := uc.GetByPvzId(ctx, pvzId)

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

func TestReceptionUseCase_CloseLastOpenOnPvz(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pvzId := "pvz-uuid-789"

	fixedTime, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	receptionOpen := entity.Reception{
		ID:       "uuid4",
		PvzID:    pvzId,
		Datetime: fixedTime,
		Status:   entity.InProgress,
	}

	receptionClose := entity.Reception{
		ID:       "uuid4",
		PvzID:    pvzId,
		Datetime: fixedTime,
		Status:   entity.Close,
	}

	tests := []struct {
		name string
		mock func(repo *mocks.MockReceptionRepo)
		res  *entity.Reception
		err  error
	}{
		{
			name: "Успешное закрытие приёмки",
			mock: func(repo *mocks.MockReceptionRepo) {
				repo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{receptionOpen}, nil)
				repo.EXPECT().
					UpdateStatus(ctx, receptionOpen.ID, entity.Close).
					Return(nil)
			},
			res: &receptionClose,
			err: nil,
		},
		{
			name: "Открытая приёмка не найдена",
			mock: func(repo *mocks.MockReceptionRepo) {
				repo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{}, nil)
			},
			res: nil,
			err: entity.ErrNoOpenReception,
		},
		{
			name: "Ошибка обновления статуса",
			mock: func(repo *mocks.MockReceptionRepo) {
				repo.EXPECT().
					GetByPvzIdWithStatus(ctx, pvzId, entity.InProgress).
					Return([]entity.Reception{receptionOpen}, nil)
				repo.EXPECT().
					UpdateStatus(ctx, receptionOpen.ID, entity.Close).
					Return(entity.ErrInternalServErr)
			},
			res: nil,
			err: entity.ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			uc, repo := receptionUseCase(t)

			localTc.mock(repo)

			res, err := uc.CloseLastOpenOnPvz(ctx, pvzId)

			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
