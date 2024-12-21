package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/injector"

	route "github.com/shigaichi/top-sites-ranking-api/internal/adapter/http"
	"github.com/shigaichi/top-sites-ranking-api/internal/adapter/http/handler"
	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := util.SetupLogger()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to set up logger when start up api server")
		return
	}

	db, err := infra.NewDb()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to create db connection when start up api server")
		return
	}

	u := injector.NewRankHistoryInteractor(db)
	h := handler.NewGetRankingImpl(u)
	ri := route.NewRouteImpl(h)
	r := ri.InitRoute()

	srv := http.Server{
		Addr:              ":3333",
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Minute,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.WithFields(log.Fields{"error": err}).Info("shutting down server")
			} else {
				log.WithFields(log.Fields{"error": err}).Error("Failed to start server")
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	var wait = 30 * time.Second

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Shutdown(ctx)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed to shutdown server")
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("Shutting down")
}
