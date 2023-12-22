package repository

import (
	"context"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type TrancoListsRepository interface {
	ExistsID(ctx context.Context, id string) (bool, error)
	Save(ctx context.Context, list model.TrancoList) error
	FindByCreatedOnLessThan(ctx context.Context, date time.Time) ([]model.TrancoList, error)
	DeleteByID(ctx context.Context, id string) error
}
