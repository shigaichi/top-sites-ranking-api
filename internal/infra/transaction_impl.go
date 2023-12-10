package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var txKey = struct{}{}

type Tx struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *Tx {
	return &Tx{db: db}
}

func (t *Tx) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, &txKey, tx)

	v, err := f(ctx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("sql execution error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("commit error: %v, rollback error: %v", err, rbErr)
		}
		return nil, err
	}
	return v, nil
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(&txKey).(*sqlx.Tx)
	return tx, ok
}
