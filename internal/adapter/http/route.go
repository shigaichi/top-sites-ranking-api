package http

import (
	"github.com/go-chi/cors"
	cfconnectingip "github.com/shigaichi/cf-connecting-ip"
	"github.com/shigaichi/top-sites-ranking-api/internal/adapter/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	mymiddleware "github.com/shigaichi/top-sites-ranking-api/internal/adapter/http/middleware"
)

type Route interface {
	InitRoute() chi.Route
}

type RouteImpl struct {
	h handler.GetRanking
}

func NewRouteImpl(h handler.GetRanking) *RouteImpl {
	return &RouteImpl{h: h}
}

func (i RouteImpl) InitRoute() chi.Router {
	router := chi.NewRouter()
	router.Use(cfconnectingip.SetRemoteAddr)
	//router.Use(middleware.Logger)
	router.Use(mymiddleware.RequestLoggerMiddleware([]string{"/status"}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/status"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           600,
	}))

	router.Route("/api/v1/rankings", func(r chi.Router) {
		r.Get("/daily", i.h.GetDailyRanking)
		r.Get("/monthly", i.h.GetMonthlyRanking)
	})

	return router
}
