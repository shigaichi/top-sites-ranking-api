package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/injector"

	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
)

func StandardWriter(date time.Time) error {
	db, err := infra.NewDb()
	if err != nil {
		return fmt.Errorf("failed to create db connection when start up service. error: %w", err)
	}
	transaction := infra.NewTransaction(db)
	u := injector.NewStandardWriteInteractor(transaction, db, 10000)

	err = u.Write(context.Background(), date)
	if err != nil {
		return fmt.Errorf("failed to write csv. error: %w", err)
	}
	return nil
}
