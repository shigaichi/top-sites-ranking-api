package infra

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type MockRankingDB struct {
	shouldError bool
}

func (m MockRankingDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	panic("no implementation")
}

func (m MockRankingDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	panic("no implementation")
}

func (m MockRankingDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}

	if query != "INSERT INTO tranco_rankings (domain_id, list_id, ranking) VALUES ($1, $2, $3),($4, $5, $6)" && query != "INSERT INTO tranco_rankings (domain_id, list_id, ranking) VALUES ($1, $2, $3)" {
		return nil, errors.New("mock error")
	}

	return nil, nil
}

func (m MockRankingDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	panic("no implementation")
}

func TestTrancoRankingsRepositoryImpl_BulkSave(t *testing.T) {
	tests := []struct {
		name    string
		args    []model.TrancoRanking
		wantErr bool
	}{
		{
			name: "successful bulk save",
			args: []model.TrancoRanking{
				{DomainID: 1, ListID: "list1", Ranking: 100},
				{DomainID: 2, ListID: "list1", Ranking: 200},
				{DomainID: 3, ListID: "list1", Ranking: 300},
			},
			wantErr: false,
		},
		{
			name: "successful single save",
			args: []model.TrancoRanking{
				{DomainID: 1, ListID: "list1", Ranking: 100},
			},
			wantErr: false,
		},
		{
			name: "failed to save",
			args: []model.TrancoRanking{
				{DomainID: 1, ListID: "list1", Ranking: 100},
				{DomainID: 2, ListID: "list1", Ranking: 200},
				{DomainID: 3, ListID: "list1", Ranking: 300},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTrancoRankingsRepositoryImpl(2, &MockRankingDB{shouldError: tt.wantErr})
			if err := r.BulkSave(context.Background(), tt.args); (err != nil) != tt.wantErr {
				t.Errorf("executeBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
