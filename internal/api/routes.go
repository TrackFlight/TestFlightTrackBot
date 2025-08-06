package api

import (
	"github.com/Laky-64/TestFlightTrackBot/internal/api/handlers"
	"github.com/Laky-64/TestFlightTrackBot/internal/api/middleware"
	"github.com/Laky-64/TestFlightTrackBot/internal/config"
	"github.com/Laky-64/TestFlightTrackBot/internal/db"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func Start(dbCtx *db.DB, cfg *config.Config) {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(chiMiddleware.RealIP)

	r.Route("/api", func(api chi.Router) {
		api.Post("/auth", handlers.AuthHandler(cfg.TelegramToken))
		api.Get("/get_config", handlers.GetConfig(cfg))

		api.Group(func(private chi.Router) {
			private.Use(middleware.JWT)

			private.Route("/users", func(users chi.Router) {
				users.Group(func(internal chi.Router) {
					internal.Use(
						middleware.AntiFlood(
							7,
							5*time.Second,
							5*time.Second,
							time.Hour,
							4*time.Minute,
						),
					)

					internal.Post("/links", handlers.AddLink(dbCtx, cfg))
					internal.Get("/links", handlers.GetLinks(dbCtx))
					internal.Delete("/links", handlers.DeleteLinks(dbCtx))
				})
				users.Group(func(internal chi.Router) {
					internal.Use(
						middleware.AntiFlood(
							10,
							5*time.Second,
							5*time.Second,
							time.Hour,
							4*time.Minute,
						),
					)

					internal.Patch("/links/{id}", handlers.EditLinkSettings(dbCtx, cfg))
				})
			})

			private.Route("/langpack", func(help chi.Router) {
				help.Use(
					middleware.AntiFlood(
						12,
						5*time.Second,
						5*time.Second,
						time.Hour,
						4*time.Minute,
					),
				)
				help.Get("/", handlers.GetLangPack(dbCtx))
			})
		})
	})

	go func() {
		if err := http.ListenAndServe(":9045", r); err != nil {
			log.Fatalf("web server error: %v", err)
		}
	}()
}
