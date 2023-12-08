package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

type TrancoDomainRepositoryImpl struct {
	db util.Crudable
}

func NewTrancoDomainRepositoryImpl(db *sqlx.DB) *TrancoDomainRepositoryImpl {
	return &TrancoDomainRepositoryImpl{db: db}
}

func (t TrancoDomainRepositoryImpl) GetIDByDomain(ctx context.Context, domain string) (int, error) {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	var id int
	query := `SELECT id FROM tranco_domains WHERE DOMAIN = $1`
	err := dao.GetContext(ctx, &id, query, domain)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error retrieving id for domain %s: %w", domain, err)
	}
	return id, nil
}

func (t TrancoDomainRepositoryImpl) Save(ctx context.Context, domain string) (int, error) {
	var dao util.Crudable
	dao, ok := GetTx(ctx)
	if !ok {
		dao = t.db
	}

	var id int
	query := `INSERT INTO tranco_domains (domain) VALUES ($1) RETURNING id`
	err := dao.QueryRowContext(ctx, query, domain).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error saving domain %s: %w", domain, err)
	}
	return id, nil
}
