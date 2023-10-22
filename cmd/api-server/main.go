package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/adapter/http/handler"
	"github.com/shigaichi/top-sites-ranking-api/internal/infra"
	"github.com/shigaichi/top-sites-ranking-api/internal/usecase"

	route "github.com/shigaichi/top-sites-ranking-api/internal/adapter/http"
	"github.com/shigaichi/top-sites-ranking-api/internal/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	if _, ok := os.LookupEnv("PROFILE"); ok {
		h := util.SetupLogger()
		defer h()
	}

	db, err := infra.NewDb()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to create db connection when start up api server")
		return
	}

	repo := infra.NewTrancoDailyRankRepositoryImpl(db)
	u := usecase.NewRankHistoryInteractor(repo)
	h := handler.NewGetRankingImpl(u)
	ri := route.NewRouteImpl(h)
	r := ri.InitRoute()

	srv := http.Server{
		Addr:    ":3333",
		Handler: r,
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
