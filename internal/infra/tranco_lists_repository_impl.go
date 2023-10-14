package infra

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

type TrancoListRepositoryImpl struct {
	db util.Crudable
}

func NewTrancoListRepositoryImpl(db *sqlx.DB) *TrancoListRepositoryImpl {
	return &TrancoListRepositoryImpl{db: db}
}

func (t TrancoListRepositoryImpl) ExistsId(ctx context.Context, id string) (bool, error) {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	var count int
	query := `SELECT COUNT(id) FROM tranco_lists WHERE id = $1`
	err := dao.GetContext(ctx, &count, query, id)
	if err != nil {
		return false, fmt.Errorf("failed tp check existence of TrancoList with ID %s: %w", id, err)
	}
	return count > 0, nil
}

func (t TrancoListRepositoryImpl) Save(ctx context.Context, list model.TrancoList) error {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	query := `INSERT INTO tranco_lists (id, created_on) VALUES ($1, $2)`
	_, err := dao.ExecContext(ctx, query, list.Id, list.CreatedOn)
	if err != nil {
		return fmt.Errorf("failed to save TrancoList with ID %s: %w", list.Id, err)
	}
	return nil
}
