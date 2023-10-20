package main

import (
	"flag"
	"os"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	if _, ok := os.LookupEnv("PROFILE"); ok {
		h := util.SetupLogger()
		defer h()
	}

	dateStr := flag.String("date", "", "Specify date in the format YYYY-MM-DD. If not specified, uses the current date.")
	flag.Parse()

	var date time.Time
	if *dateStr == "" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", *dateStr)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "date": date}).Error("Error parsing date")
			return
		}
	}

	err := internal.StandardWriter(date)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "date": date}).Error("Error in StandardWriter")
		return
	}
}
