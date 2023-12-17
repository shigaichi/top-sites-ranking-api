//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
	"github.com/shigaichi/top-sites-ranking-api/internal/usecase"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

func NewRankHistoryInteractor(db util.Crudable) *usecase.RankHistoryInteractor {
	wire.Build(
		usecase.NewRankHistoryInteractor,
		infra.NewTrancoDailyRankRepositoryImpl,
		wire.Bind(new(repository.TrancoDailyRankRepository), new(*infra.TrancoDailyRankRepositoryImpl)),
	)
	return nil
}

func NewStandardWriteInteractor(transaction repository.Transaction, db *sqlx.DB, batchSize int) *usecase.StandardWriteInteractor {
	wire.Build(
		usecase.NewStandardWriteInteractor,
		infra.NewTrancoAPIImpl,
		wire.Bind(new(repository.TrancoAPIRepository), new(*infra.TrancoAPIImpl)),
		infra.NewTrancoListRepositoryImpl,
		wire.Bind(new(repository.TrancoListsRepository), new(*infra.TrancoListRepositoryImpl)),
		infra.NewTrancoCsvImpl,
		wire.Bind(new(repository.TrancoCsvRepository), new(*infra.TrancoCsvImpl)),
		infra.NewTrancoDomainRepositoryImpl,
		wire.Bind(new(repository.TrancoDomainsRepository), new(*infra.TrancoDomainRepositoryImpl)),
		infra.NewTrancoRankingsRepositoryImpl,
		wire.Bind(new(repository.TrancoRankingsRepository), new(*infra.TrancoRankingsRepositoryImpl)),
		wire.Bind(new(util.Crudable), new(*sqlx.DB)),
	)
	return nil
}
