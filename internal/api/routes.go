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

		api.Group(func(private chi.Router) {
			private.Use(middleware.JWT)
			private.Use(
				middleware.AntiFlood(
					7,
					5*time.Second,
					5*time.Second,
					time.Hour,
					4*time.Minute,
				),
			)

			private.Route("/users", func(users chi.Router) {
				users.Get("/links", handlers.GetLinks(dbCtx))
				users.Post("/links", handlers.AddLink(dbCtx))
				users.Delete("/links/{linkID}", handlers.DeleteLink(dbCtx))
			})

			private.Route("/langpack", func(help chi.Router) {
				help.Get("/strings", handlers.GetStrings(dbCtx))
			})
		})
	})

	go func() {
		if err := http.ListenAndServe(":9045", r); err != nil {
			log.Fatalf("web server error: %v", err)
		}
	}()
}
