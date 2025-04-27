package usecase

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/shigaichi/tranco"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
	log "github.com/sirupsen/logrus"
)

type WriteUseCase interface {
	Write(ctx context.Context, date time.Time) error
}

type StandardWriteInteractor struct {
	api         repository.TrancoAPIRepository
	list        repository.TrancoListsRepository
	csv         repository.TrancoCsvRepository
	transaction repository.Transaction
	domain      repository.TrancoDomainsRepository
	ranking     repository.TrancoRankingsRepository
}

func NewStandardWriteInteractor(api repository.TrancoAPIRepository, list repository.TrancoListsRepository, csv repository.TrancoCsvRepository, transaction repository.Transaction, domain repository.TrancoDomainsRepository, ranking repository.TrancoRankingsRepository) *StandardWriteInteractor {
	return &StandardWriteInteractor{api: api, list: list, csv: csv, transaction: transaction, domain: domain, ranking: ranking}
}

func (i StandardWriteInteractor) Write(ctx context.Context, date time.Time) error {
	const maxRetries = 3
	const retryInterval = 100 * time.Millisecond

	// retry get tranco id until success
	var metadata tranco.ListMetadata
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		var err error
		metadata, err = i.api.GetIDByDate(date)
		if err == nil {
			break
		}
		lastErr = err
		time.Sleep(retryInterval)
	}

	if lastErr != nil {
		return fmt.Errorf("failed to get tranco list id by date after %d retries: %w", maxRetries, lastErr)
	}

	savedListID, err := i.list.ExistsID(ctx, metadata.ListID)
	if err != nil {
		return fmt.Errorf("failed to check list id is already exist or not in writing standard tranco list error: %w", err)
	}

	if savedListID {
		log.WithFields(log.Fields{"list_id": metadata.ListID, "date": date}).Info("list id already exists in writing standard tranco list")
		return nil
	} else {
		log.WithFields(log.Fields{"list_id": metadata.ListID, "date": date}).Info("list id does not exist and write standard tranco list")
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
		err = i.list.Save(ctx, model.TrancoList{ID: metadata.ListID, CreatedOn: metadata.CreatedOn})
		if err != nil {
			return nil, fmt.Errorf("failed to save tranco list with id %s error: %w", metadata.ListID, err)
		}

		var l []model.TrancoRanking
		for _, ranking := range rankings {
			var domainID int
			domainID, err = i.domain.GetIDByDomain(ctx, ranking.Domain)
			if err != nil {
				return nil, fmt.Errorf("failed to save list id in writing standard tranco list error: %w", err)
			}

			if domainID == 0 {
				domainID, err = i.domain.Save(ctx, ranking.Domain)
				if err != nil {
					return nil, fmt.Errorf("failed to save domain in writing standard tranco list error: %w", err)
				}
			}

			l = append(l, model.TrancoRanking{DomainID: domainID, ListID: metadata.ListID, Ranking: ranking.Rank})
		}

		err := i.ranking.BulkSave(ctx, l)
		if err != nil {
			return nil, fmt.Errorf("failed to bulk save %d rankings in writing standard tranco list error: %w", len(l), err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: %w", err)
	}

	return nil
}
