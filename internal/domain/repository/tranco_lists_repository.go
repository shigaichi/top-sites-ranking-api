package repository

import (
	"context"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type TrancoListsRepository interface {
	ExistsID(ctx context.Context, id string) (bool, error)
	Save(ctx context.Context, list model.TrancoList) error
}
