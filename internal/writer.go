package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
	"github.com/shigaichi/top-sites-ranking-api/internal/usecase"
)

func StandardWriter(date time.Time) error {
	db, err := infra.NewDb()
	if err != nil {
		return fmt.Errorf("failed to create db connection when start up service. error: %w", err)
	}
	api := infra.NewTrancoAPIImpl()
	csv := infra.NewTrancoCsvImpl()
	lists := infra.NewTrancoListRepositoryImpl(db)
	domain := infra.NewTrancoDomainRepositoryImpl(db)
	rankings := infra.NewTrancoRankingsRepositoryImpl(10000, db)
	transaction := infra.NewTransaction(db)

	u := usecase.NewStandardWriteInteractor(api, lists, csv, transaction, domain, rankings)
	err = u.Write(context.Background(), date)
	if err != nil {
		return fmt.Errorf("failed to write csv. error: %w", err)
	}
	return nil
}
