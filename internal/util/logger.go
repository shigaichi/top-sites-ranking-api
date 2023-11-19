package util

import (
	"os"

	"github.com/highlight/highlight/sdk/highlight-go"
	hlog "github.com/highlight/highlight/sdk/highlight-go/log"
	log "github.com/sirupsen/logrus"
)

func SetupLogger() error {
	p := GetEnvWithDefault("LOG_LEVEL", "info")

	level, err := log.ParseLevel(p)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	return nil
}

func SetupHighlight() func() {
	value, exists := os.LookupEnv("HIGHLIGHT_ID")

	if exists {
		// set up the highlight SDK
		highlight.SetProjectID(value)
		highlight.Start(
			highlight.WithServiceName(GetEnvWithDefault("SERVICE_NAME", "top-sites-ranking-api")),
			highlight.WithServiceVersion(GetEnvWithDefault("VERSION", "dev")),
		)
		// setup highlight logrus hook
		hlog.Init()
		return highlight.Stop
	}

	return nil
}
