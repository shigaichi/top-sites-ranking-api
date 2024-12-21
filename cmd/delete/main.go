package main

import (
	"context"
	"flag"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
	"github.com/shigaichi/top-sites-ranking-api/internal/injector"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := util.SetupLogger()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to set up logger when starting delete")
	}

	var sinceFlag int
	flag.IntVar(&sinceFlag, "since", 100, "Number of days back to process")

	flag.Parse()

	if sinceFlag < 0 {
		log.WithFields(log.Fields{"since": sinceFlag}).Fatal("error: 'since' flag cannot be negative")
	}

	db, err := infra.NewDb()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to create db connection when start up service.")
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.WithFields(log.Fields{"since": sinceFlag}).Fatal("failed to close db")
		}
	}(db)

	usecase := injector.NewDeleteInteractor(db, 1)
	err = usecase.Delete(context.Background(), time.Duration(sinceFlag)*24*time.Hour)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "since": sinceFlag}).Fatal("failed to delete old records")
	}
}
