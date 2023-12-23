package infra

import (
	"context"
	"fmt"
	"time"

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

func (t TrancoListRepositoryImpl) ExistsID(ctx context.Context, id string) (bool, error) {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	var count int
	query := `SELECT COUNT(id) FROM tranco_lists WHERE id = $1`
	err := dao.GetContext(ctx, &count, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of TrancoList with ID %s: %w", id, err)
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
	_, err := dao.ExecContext(ctx, query, list.ID, list.CreatedOn)
	if err != nil {
		return fmt.Errorf("failed to save TrancoList with ID %s: %w", list.ID, err)
	}
	return nil
}

func (t TrancoListRepositoryImpl) FindByCreatedOnLessThan(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	var lists []model.TrancoList
	query := "SELECT id, created_on FROM tranco_lists WHERE created_on < $1"
	err := dao.SelectContext(ctx, &lists, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to find by created on less than %s: %w", date, err)
	}

	return lists, nil
}

func (t TrancoListRepositoryImpl) DeleteByID(ctx context.Context, id string) error {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	query := "DELETE FROM tranco_lists WHERE id = $1"
	_, err := dao.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete by id %s: %w", id, err)
	}
	return nil
}
