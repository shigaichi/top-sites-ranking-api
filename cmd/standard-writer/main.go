package main

import (
	"flag"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	dateStr := flag.String("date", "", "Specify date in the format YYYY-MM-DD. If not specified, uses the current date.")
	flag.Parse()

	var date time.Time
	if *dateStr == "" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", *dateStr)
		if err != nil {
			log.Fatalf("Error parsing date: %v", err)
			return
		}
	}

	err := internal.StandardWriter(date)
	if err != nil {
		log.Fatalf("Error in StandardWriter: %v", err)
	}
}
