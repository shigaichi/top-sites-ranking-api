package repository

import (
	"context"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type TrancoRankingsRepository interface {
	Save(ctx context.Context, ranking model.TrancoRanking) error
	BulkSave(ctx context.Context, rankings []model.TrancoRanking) error
}
