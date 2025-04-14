package pvz

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pvz-service/internal/entity"
	mocks "pvz-service/internal/repo"
	"testing"
	"time"
)

var errInternalServErr = errors.New("internal server error")

func pvzUseCase(t *testing.T) (*UseCase, *mocks.MockPvzRepo) {
	t.Helper()

	pvzRepo := mocks.NewMockPvzRepo(t)
	useCase := New(pvzRepo)

	return useCase, pvzRepo
}

func TestPvzUseCase_Create(t *testing.T) {
	t.Parallel()

	uc, repo := pvzUseCase(t)

	fixedDate, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name  string
		input entity.Pvz
		mock  func()
		res   *entity.Pvz
		err   error
	}{
		{
			name:  "Успешно создан ПВЗ",
			input: entity.Pvz{ID: "uuid4", RegistrationDate: fixedDate, City: "Moscow"},
			mock: func() {
				pvzFake := entity.Pvz{
					ID:               "uuid4",
					RegistrationDate: fixedDate,
					City:             "Moscow",
				}
				repo.EXPECT().Create(ctx, pvzFake).Return(&pvzFake, nil)
			},
			res: &entity.Pvz{ID: "uuid4", RegistrationDate: fixedDate, City: "Moscow"},
			err: nil,
		},
		{
			name:  "Репо возвращает ошибку",
			input: entity.Pvz{ID: "uuid-error", RegistrationDate: fixedDate, City: "ErrorCity"},
			mock: func() {
				pvzFake := entity.Pvz{
					ID:               "uuid-error",
					RegistrationDate: fixedDate,
					City:             "ErrorCity",
				}
				repo.EXPECT().Create(ctx, pvzFake).Return(nil, errInternalServErr)
			},
			res: nil,
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, err := uc.Create(ctx, localTc.input)

			assert.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}

type listTestInput struct {
	startDate *time.Time
	endDate   *time.Time
	page      int
	limit     int
}

func TestPvzUseCase_ListByReceptionDate(t *testing.T) {
	t.Parallel()

	uc, repo := pvzUseCase(t)

	ctx := context.Background()

	startDate1, err := time.Parse(time.DateOnly, "2002-06-21")
	require.NoError(t, err)
	endDate1, err := time.Parse(time.DateOnly, "2002-06-28")
	require.NoError(t, err)

	startDate2, err := time.Parse(time.DateOnly, "2003-06-21")
	require.NoError(t, err)
	endDate2, err := time.Parse(time.DateOnly, "2003-06-21")
	require.NoError(t, err)

	pvz1 := entity.Pvz{
		ID:               "1",
		RegistrationDate: startDate1,
		City:             "City1",
	}
	pvz2 := entity.Pvz{
		ID:               "2",
		RegistrationDate: endDate1,
		City:             "City2",
	}

	tests := []struct {
		name  string
		input listTestInput
		mock  func()
		res   []entity.Pvz
		err   error
	}{
		{
			name: "Успешно с датами приемки",
			input: listTestInput{
				startDate: &startDate1,
				endDate:   &endDate1,
				page:      1,
				limit:     10,
			},
			mock: func() {
				repo.EXPECT().GetByReceptionDate(ctx, &startDate1, &endDate1, 10, 0).
					Return([]entity.Pvz{pvz1, pvz2}, nil)
			},
			res: []entity.Pvz{pvz1, pvz2},
			err: nil,
		},
		{
			name: "Репо возвращает ошибку",
			input: listTestInput{
				startDate: &startDate2,
				endDate:   &endDate2,
				page:      1,
				limit:     10,
			},
			mock: func() {
				repo.EXPECT().GetByReceptionDate(ctx, &startDate2, &endDate2, 10, 0).
					Return(nil, errInternalServErr)
			},
			res: nil,
			err: errInternalServErr,
		},
		{
			name: "Успешно без дат приемки",
			input: listTestInput{
				startDate: nil,
				endDate:   nil,
				page:      1,
				limit:     10,
			},
			mock: func() {
				repo.EXPECT().GetByReceptionDate(ctx, (*time.Time)(nil), (*time.Time)(nil), 10, 0).
					Return([]entity.Pvz{pvz1}, nil)

			},
			res: []entity.Pvz{pvz1},
			err: nil,
		},
	}
	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()
			res, err := uc.ListByReceptionDate(
				ctx,
				localTc.input.startDate,
				localTc.input.endDate,
				localTc.input.page,
				localTc.input.limit,
			)
			require.Equal(t, localTc.res, res)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
