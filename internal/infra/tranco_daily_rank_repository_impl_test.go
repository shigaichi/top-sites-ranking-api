package infra

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type MockDB struct {
	shouldError   bool
	returnNoRanks bool
}

func (m MockDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if m.shouldError {
		return errors.New("mock error")
	}
	ranks := dest.(*[]model.DailyRank)
	if m.returnNoRanks {
		*ranks = []model.DailyRank{}
		return nil
	}
	*ranks = []model.DailyRank{{Rank: 1, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}}
	return nil
}

func (m MockDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	panic("no implementation")
}

func (m MockDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	panic("no implementation")
}

func (m MockDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	panic("no implementation")
}

func TestGetDailyRanksByDateRange(t *testing.T) {
	tests := []struct {
		name          string
		shouldError   bool
		returnNoRanks bool
		wantError     bool
		wantRanks     []model.DailyRank
	}{
		{
			name:        "successful fetch",
			shouldError: false,
			wantError:   false,
			wantRanks:   []model.DailyRank{{Rank: 1, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)}},
		},
		{
			name:        "DB error",
			shouldError: true,
			wantError:   true,
		},
		{
			name:          "0 fetch",
			shouldError:   false,
			returnNoRanks: true,
			wantError:     false,
			wantRanks:     []model.DailyRank{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := TrancoDailyRankRepositoryImpl{db: &MockDB{shouldError: tt.shouldError, returnNoRanks: tt.returnNoRanks}}
			ranks, err := repo.GetDailyRanksByDateRange(context.Background(), "example.com", time.Now(), time.Now())
			if tt.wantError {
				if err == nil {
					t.Errorf("expected an error, got nil")
					return
				}
				return
			}
			if err != nil {
				t.Errorf("didn't expect an error, got %v", err)
				return
			}

			if diff := cmp.Diff(ranks, tt.wantRanks); diff != "" {
				t.Errorf("result is mimatch:\n%s", diff)
				return
			}
		})
	}
}
