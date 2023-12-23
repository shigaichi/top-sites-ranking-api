package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
)

type DeleteUseCase interface {
	Delete(ctx context.Context, duration time.Duration) error
}

type DeleteInteractor struct {
	list    repository.TrancoListsRepository
	ranking repository.TrancoRankingsRepository
}

func NewDeleteInteractor(list repository.TrancoListsRepository, ranking repository.TrancoRankingsRepository) *DeleteInteractor {
	return &DeleteInteractor{list: list, ranking: ranking}
}

// Delete removes TrancoLists created before a given duration and their associated TrancoRankings if CreatedOn is not the end of the month.
func (d DeleteInteractor) Delete(ctx context.Context, duration time.Duration) error {
	cutoffDate := time.Now().Add(-duration)

	lists, err := d.list.FindByCreatedOnLessThan(ctx, cutoffDate)
	if err != nil {
		return fmt.Errorf("error finding tranco lists: %w", err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(lists))

	for _, list := range lists {
		if !isEndOfMonth(list.CreatedOn) {
			wg.Add(1)
			go func(list model.TrancoList) {
				defer wg.Done()

				log.WithFields(log.Fields{"listID": list.ID, "CreatedOn": list.CreatedOn}).Info("delete list and ranking")

				if err := d.deleteListAndRankings(ctx, list.ID); err != nil {
					errs <- fmt.Errorf("error deleting list and rankings for list ID %s: %w", list.ID, err)
				}
			}(list)
		}
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DeleteInteractor) deleteListAndRankings(ctx context.Context, listID string) error {
	if err := d.ranking.DeleteByListID(ctx, listID); err != nil {
		return fmt.Errorf("error deleting tranco rankings by list ID %s: %w", listID, err)
	}

	if err := d.list.DeleteByID(ctx, listID); err != nil {
		return fmt.Errorf("error deleting tranco list by ID %s: %w", listID, err)
	}

	return nil
}

func isEndOfMonth(date time.Time) bool {
	endOfMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location()).Add(-24 * time.Hour)
	return date.Day() == endOfMonth.Day()
}
