package infra

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

type TrancoRankingsRepositoryImpl struct {
	db util.Crudable
}

func NewTrancoRankingsRepositoryImpl(db *sqlx.DB) *TrancoRankingsRepositoryImpl {
	return &TrancoRankingsRepositoryImpl{db: db}
}

func (t TrancoRankingsRepositoryImpl) Save(ctx context.Context, ranking model.TrancoRanking) error {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	query := `INSERT INTO tranco_rankings (domain_id, list_id, ranking) VALUES ($1, $2, $3)`

	_, err := dao.ExecContext(ctx, query, ranking.DomainId, ranking.ListId, ranking.Ranking)
	if err != nil {
		return fmt.Errorf("error saving tranco ranking: %w", err)
	}
	return nil
}
