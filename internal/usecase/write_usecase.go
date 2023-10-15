package usecase

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
	log "github.com/sirupsen/logrus"
)

type WriteUseCase interface {
	Write(ctx context.Context, date time.Time) error
}

type StandardWriteInteractor struct {
	api         repository.TrancoApiRepository
	csv         repository.TrancoCsvRepository
	list        repository.TrancoListsRepository
	domain      repository.TrancoDomainsRepository
	ranking     repository.TrancoRankingsRepository
	transaction repository.Transaction
}

func NewStandardWriteInteractor(api repository.TrancoApiRepository, csv repository.TrancoCsvRepository, list repository.TrancoListsRepository, domain repository.TrancoDomainsRepository, ranking repository.TrancoRankingsRepository, transaction repository.Transaction) *StandardWriteInteractor {
	return &StandardWriteInteractor{api: api, csv: csv, list: list, domain: domain, ranking: ranking, transaction: transaction}
}

func (i StandardWriteInteractor) Write(ctx context.Context, date time.Time) error {
	metadata, err := i.api.GetIdByDate(date)
	if err != nil {
		return fmt.Errorf("failed to get tranco list id by date in writing standard tranco list error: %w", err)
	}

	savedListId, err := i.list.ExistsId(ctx, metadata.ListId)
	if err != nil {
		return fmt.Errorf("failed to check list id is already exist or not in writing standard tranco list error: %w", err)
	}

	if savedListId {
		log.Infof("list id %s (date: %s) alread exists in writing standard tranco list", metadata.ListId, date.String())
		return nil
	} else {
		log.Infof("list id %s (date: %s) does not exist and write standard tranco list", metadata.ListId, date.String())
	}

	parse, err := url.Parse(metadata.Download)
	if err != nil {
		return fmt.Errorf("failed to parse csv url in writing standard tranco list error: %w", err)
	}
	rankings, err := i.csv.Get(*parse)
	if err != nil {
		return fmt.Errorf("failed to get csv in writing standard tranco list error: %w", err)
	}

	_, err = i.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {

		err = i.list.Save(ctx, model.TrancoList{Id: metadata.ListId, CreatedOn: metadata.CreatedOn})
		if err != nil {
			return nil, fmt.Errorf("failed to save list id is already exist or not in writing standard tranco list error: %w", err)
		}

		for _, ranking := range rankings {
			var domainId int
			domainId, err = i.domain.GetIdByDomain(ctx, ranking.Domain)
			if err != nil {
				return nil, fmt.Errorf("failed to save list id in writing standard tranco list error: %w", err)
			}

			if domainId == 0 {
				domainId, err = i.domain.Save(ctx, ranking.Domain)
				if err != nil {
					return nil, fmt.Errorf("failed to save domain in writing standard tranco list error: %w", err)
				}
			}

			err := i.ranking.Save(ctx, model.TrancoRanking{DomainId: domainId, ListId: metadata.ListId, Ranking: ranking.Rank})
			if err != nil {
				return nil, fmt.Errorf("failed to save ranking in writing standard tranco list error: %w", err)
			}
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: %w", err)
	}

	return nil
}