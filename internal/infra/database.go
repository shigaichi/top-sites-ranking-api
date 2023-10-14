package infra

import (
	"fmt"

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
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database in db configuration. error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database in db configuration. error: %w", err)
	}
	return db, nil
}
