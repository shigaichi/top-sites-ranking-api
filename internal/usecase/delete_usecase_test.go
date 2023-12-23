package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
)

type MockTrancoListsRepositoryForDelete struct {
	infra.TrancoListRepositoryImpl
	MockFindByCreatedOnLessThan func(ctx context.Context, date time.Time) ([]model.TrancoList, error)
	MockDeleteById              func(ctx context.Context, id string) error
}

func (m MockTrancoListsRepositoryForDelete) FindByCreatedOnLessThan(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
	return m.MockFindByCreatedOnLessThan(ctx, date)
}

func (m MockTrancoListsRepositoryForDelete) DeleteByID(ctx context.Context, id string) error {
	return m.MockDeleteById(ctx, id)
}

type MockTrancoRankingsRepositoryForDelete struct {
	infra.TrancoRankingsRepositoryImpl
	MockDeleteByListID func(ctx context.Context, listID string) error
}

func (m MockTrancoRankingsRepositoryForDelete) DeleteByListID(ctx context.Context, listID string) error {
	return m.MockDeleteByListID(ctx, listID)
}

func TestDeleteInteractor_Delete(t *testing.T) {
	tests := []struct {
		name      string
		duration  time.Duration
		setupMock func(*MockTrancoListsRepositoryForDelete, *MockTrancoRankingsRepositoryForDelete)
		wantErr   bool
	}{
		{
			name:     "Successful deletion",
			duration: 24 * time.Hour,
			setupMock: func(mList *MockTrancoListsRepositoryForDelete, mRankings *MockTrancoRankingsRepositoryForDelete) {
				mList.MockFindByCreatedOnLessThan = func(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
					return []model.TrancoList{{ID: "Q94V4", CreatedOn: time.Date(2023, 1, 2, 3, 4, 5, 6, time.Local)}}, nil
				}
				mList.MockDeleteById = func(ctx context.Context, id string) error {
					return nil
				}
				mRankings.MockDeleteByListID = func(ctx context.Context, listID string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:     "No delete",
			duration: 24 * time.Hour,
			setupMock: func(mList *MockTrancoListsRepositoryForDelete, mRankings *MockTrancoRankingsRepositoryForDelete) {
				mList.MockFindByCreatedOnLessThan = func(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
					return []model.TrancoList{}, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "Only end of month",
			duration: 24 * time.Hour,
			setupMock: func(mList *MockTrancoListsRepositoryForDelete, mRankings *MockTrancoRankingsRepositoryForDelete) {
				mList.MockFindByCreatedOnLessThan = func(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
					return []model.TrancoList{{ID: "Q94V4", CreatedOn: time.Date(2023, 1, 31, 0, 0, 0, 0, time.Local)}}, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "Error in find list",
			duration: 24 * time.Hour,
			setupMock: func(mList *MockTrancoListsRepositoryForDelete, mRankings *MockTrancoRankingsRepositoryForDelete) {
				mList.MockFindByCreatedOnLessThan = func(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
					return nil, errors.New("mock error")
				}
			},
			wantErr: true,
		},
		{
			name:     "Error in delete ranks",
			duration: 24 * time.Hour,
			setupMock: func(mList *MockTrancoListsRepositoryForDelete, mRankings *MockTrancoRankingsRepositoryForDelete) {
				mList.MockFindByCreatedOnLessThan = func(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
					return []model.TrancoList{{ID: "Q94V4", CreatedOn: time.Date(2023, 1, 2, 3, 4, 5, 6, time.Local)}}, nil
				}
				mRankings.MockDeleteByListID = func(ctx context.Context, listID string) error {
					return errors.New("mock error")
				}
			},
			wantErr: true,
		},
		{
			name:     "Error in delete lists",
			duration: 24 * time.Hour,
			setupMock: func(mList *MockTrancoListsRepositoryForDelete, mRankings *MockTrancoRankingsRepositoryForDelete) {
				mList.MockFindByCreatedOnLessThan = func(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
					return []model.TrancoList{{ID: "Q94V4", CreatedOn: time.Date(2023, 1, 2, 3, 4, 5, 6, time.Local)}}, nil
				}
				mRankings.MockDeleteByListID = func(ctx context.Context, listID string) error {
					return nil
				}
				mList.MockDeleteById = func(ctx context.Context, id string) error {
					return errors.New("mock error")
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockListRepo := &MockTrancoListsRepositoryForDelete{}
			mockRankingRepo := &MockTrancoRankingsRepositoryForDelete{}
			tt.setupMock(mockListRepo, mockRankingRepo)

			d := DeleteInteractor{
				list:    mockListRepo,
				ranking: mockRankingRepo,
			}

			err := d.Delete(context.Background(), tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
