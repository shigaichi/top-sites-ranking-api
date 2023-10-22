package repository

import (
	"context"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type TrancoDailyRankRepository interface {
	GetDailyRanksByDateRange(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error)
}
