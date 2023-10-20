package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	route "github.com/shigaichi/top-sites-ranking-api/internal/adapter/http"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	if _, ok := os.LookupEnv("PROFILE"); ok {
		h := util.SetupLogger()
		defer h()
	}

	r := route.InitRoute()

	srv := http.Server{
		Addr:    ":3333",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to start server")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	var wait time.Duration

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("Shutting down")
}
