package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
)

type RankHistoryUseCase interface {
	GetDailyRanking(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error)
	GetMonthlyRanking(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error)
}

type RankHistoryInteractor struct {
	repo repository.TrancoDailyRankRepository
}

func NewRankHistoryInteractor(repo repository.TrancoDailyRankRepository) *RankHistoryInteractor {
	return &RankHistoryInteractor{repo: repo}
}

func (r RankHistoryInteractor) GetDailyRanking(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error) {
	ranks, err := r.repo.GetDailyRanksByDateRange(ctx, domain, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily ranks: %w", err)
	}
	return ranks, nil
}

func (r RankHistoryInteractor) GetMonthlyRanking(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error) {
	ranks, err := r.repo.GetDailyRanksByDateRange(ctx, domain, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily ranks: %w", err)
	}

	var monthlyRanks []model.DailyRank
	for _, rank := range ranks {
		if isLastDayOfMonth(rank.Date) {
			monthlyRanks = append(monthlyRanks, rank)
		}
	}

	return monthlyRanks, nil
}

// Helper function to check if the given date is the last day of its month
func isLastDayOfMonth(t time.Time) bool {
	nextDay := t.AddDate(0, 0, 1)
	return t.Month() != nextDay.Month()
}
