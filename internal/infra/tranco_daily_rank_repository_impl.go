package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

type TrancoDailyRankRepositoryImpl struct {
	db util.Crudable
}

func NewTrancoDailyRankRepositoryImpl(db util.Crudable) *TrancoDailyRankRepositoryImpl {
	return &TrancoDailyRankRepositoryImpl{db: db}
}

func (t TrancoDailyRankRepositoryImpl) GetDailyRanksByDateRange(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error) {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	var ranks []model.DailyRank

	query := `
SELECT tr.ranking AS Rank, tl.created_on AS Date
FROM tranco_rankings tr
         INNER JOIN tranco_domains td ON tr.domain_id = td.id
         INNER JOIN public.tranco_lists tl ON tr.list_id = tl.id
WHERE td.domain = $1
  AND tl.created_on BETWEEN $2 AND $3
  ORDER BY Date DESC
`

	args := []interface{}{domain, start, end}

	if err := dao.SelectContext(ctx, &ranks, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch daily ranks: %w", err)
	}

	return ranks, nil
}
