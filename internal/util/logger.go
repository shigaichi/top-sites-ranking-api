package util

import (
	log "github.com/sirupsen/logrus"
)

func SetupLogger() error {
	p := GetEnvWithDefault("LOG_LEVEL", "info")

	level, err := log.ParseLevel(p)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	log.SetFormatter(&log.JSONFormatter{})

	return nil
}
