package infra

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

type TrancoRankingsRepositoryImpl struct {
	batchSize int
	db        util.Crudable
}

func NewTrancoRankingsRepositoryImpl(batchSize int, db util.Crudable) *TrancoRankingsRepositoryImpl {
	return &TrancoRankingsRepositoryImpl{batchSize: batchSize, db: db}
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

func (t TrancoRankingsRepositoryImpl) BulkSave(ctx context.Context, rankings []model.TrancoRanking) error {
	for i := 0; i < len(rankings); i += t.batchSize {
		end := i + t.batchSize
		if end > len(rankings) {
			end = len(rankings)
		}
		if err := t.executeBatch(ctx, rankings[i:end]); err != nil {
			return fmt.Errorf("error saving tranco ranking batch: %w", err)
		}
	}
	return nil
}

func (t TrancoRankingsRepositoryImpl) executeBatch(ctx context.Context, rankings []model.TrancoRanking) error {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	valueStrings := make([]string, 0, len(rankings))
	valueArgs := make([]any, 0, len(rankings)*3)
	for i, ranking := range rankings {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, ranking.DomainId, ranking.ListId, ranking.Ranking)
	}

	query := fmt.Sprintf("INSERT INTO tranco_rankings (domain_id, list_id, ranking) VALUES %s", strings.Join(valueStrings, ","))

	log.WithFields(log.Fields{"rank_count": len(rankings)}).Debug("rank data saved")

	_, err := dao.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("error executing batch insert: %w", err)
	}
	return nil
}
