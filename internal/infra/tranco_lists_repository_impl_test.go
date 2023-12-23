package infra

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/google/go-cmp/cmp"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type MockListDB struct {
	sqlx.DB
	MockExecContext   func(ctx context.Context, query string, args ...any) (sql.Result, error)
	MockSelectContext func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (m MockListDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.MockExecContext(ctx, query, args)
}

func (m MockListDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.MockSelectContext(ctx, dest, query, args)
}

func TestTrancoListRepositoryImpl_FindByCreatedOnLessThan(t1 *testing.T) {
	type args struct {
		ctx  context.Context
		date time.Time
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(db *MockListDB)
		want      []model.TrancoList
		wantErr   bool
	}{
		{
			name: "Valid date with no errors",
			args: args{ctx: context.Background(), date: time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC)},
			setupMock: func(db *MockListDB) {
				db.MockSelectContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					lists := dest.(*[]model.TrancoList)
					*lists = []model.TrancoList{
						{ID: "1", CreatedOn: time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC)},
						{ID: "2", CreatedOn: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)},
					}
					return nil
				}
			},
			want: []model.TrancoList{
				{ID: "1", CreatedOn: time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC)},
				{ID: "2", CreatedOn: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
		},
		{
			name: "Valid date with error",
			args: args{ctx: context.Background(), date: time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC)},
			setupMock: func(db *MockListDB) {
				db.MockSelectContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					return errors.New("mock error")
				}
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			db := MockListDB{}
			tt.setupMock(&db)
			t := TrancoListRepositoryImpl{
				db: &db,
			}
			got, err := t.FindByCreatedOnLessThan(tt.args.ctx, tt.args.date)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindByCreatedOnLessThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t1.Errorf("result is mimatch:\n%s", diff)
			}
		})
	}
}

func TestTrancoListRepositoryImpl_DeleteById(t1 *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(db *MockListDB)
		wantErr   bool
	}{
		{
			name: "Valid delete with no errors",
			args: args{ctx: context.Background(), id: "1"},
			setupMock: func(db *MockListDB) {
				db.MockExecContext = func(ctx context.Context, query string, args ...any) (sql.Result, error) {
					return nil, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Delete with error",
			args: args{ctx: context.Background(), id: "invalid-id"},
			setupMock: func(db *MockListDB) {
				db.MockExecContext = func(ctx context.Context, query string, args ...any) (sql.Result, error) {
					return nil, errors.New("mock test")
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			db := MockListDB{}
			tt.setupMock(&db)
			t := TrancoListRepositoryImpl{
				db: &db,
			}
			if err := t.DeleteByID(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t1.Errorf("DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
