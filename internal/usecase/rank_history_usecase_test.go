package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

// モックリポジトリ
type mockRepo struct {
	data []model.DailyRank
	err  error
}

func (m *mockRepo) GetDailyRanksByDateRange(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error) {
	return m.data, m.err
}

func TestRankHistoryInteractor_GetDailyRanking(t *testing.T) {
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		repoData []model.DailyRank
		repoErr  error
		expected []model.DailyRank
		err      error
	}{
		{
			name: "Successful case",
			repoData: []model.DailyRank{
				{Rank: 1, Date: endTime},
			},
			expected: []model.DailyRank{
				{Rank: 1, Date: endTime},
			},
			err: nil,
		},
		{
			name:     "Empty data from repository",
			repoData: []model.DailyRank{},
			expected: []model.DailyRank{},
		},
		{
			name:    "Repository error",
			repoErr: errors.New("some error"),
			err:     errors.New("failed to get daily ranks: some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactor := RankHistoryInteractor{
				repo: &mockRepo{data: tt.repoData, err: tt.repoErr},
			}

			got, err := interactor.GetDailyRanking(context.TODO(), "example.com", startTime, endTime)

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}

			if tt.err != nil {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.err.Error() {
					t.Errorf("expected error: %v, got: %v", tt.err, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRankHistoryInteractor_GetMonthlyRanking(t *testing.T) {
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		repoData []model.DailyRank
		repoErr  error
		expected []model.DailyRank
		err      error
	}{
		{
			name: "Successful case",
			repoData: []model.DailyRank{
				{Rank: 1, Date: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)},
				{Rank: 2, Date: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC)},
				{Rank: 3, Date: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)},
			},
			expected: []model.DailyRank{
				{Rank: 1, Date: time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)},
				{Rank: 3, Date: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:     "Empty data from repository",
			repoData: []model.DailyRank{},
			expected: nil,
		},
		{
			name:    "Repository error",
			repoErr: errors.New("some error"),
			err:     errors.New("failed to get daily ranks: some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RankHistoryInteractor{
				repo: &mockRepo{data: tt.repoData, err: tt.repoErr},
			}

			result, err := r.GetMonthlyRanking(context.Background(), "testdomain.com", startTime, endTime)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
			if (err != nil || tt.err != nil) && err.Error() != tt.err.Error() {
				t.Errorf("expected error %v, but got %v", tt.err, err)
			}
		})
	}
}
