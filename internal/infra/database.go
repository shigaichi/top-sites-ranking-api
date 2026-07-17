package infra

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
)

func NewDb() (*sqlx.DB, error) {
	user := util.GetEnvWithDefault("DB_USER", "user")
	port := util.GetEnvWithDefault("DB_PORT", "5432")
	host := util.GetEnvWithDefault("DB_HOST", "localhost")
	password := util.GetEnvWithDefault("DB_PASSWORD", "password")
	dbname := util.GetEnvWithDefault("DB_NAME", "ranking")
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	const maxRetries = 5
	const retryInterval = 2 * time.Second

	var db *sqlx.DB
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		var err error
		db, err = sqlx.Connect("pgx", dataSourceName)
		if err == nil {
			return db, nil
		}
		lastErr = err
		log.WithFields(log.Fields{"attempt": attempt + 1, "max_retries": maxRetries}).Warn("failed to connect to database, retrying...")
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database in db configuration after %d retries. error: %w", maxRetries, lastErr)
}
