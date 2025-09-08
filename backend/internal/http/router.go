package http

import (
	"my-app-backend/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", handlers.Health)
	r.Get("/api/quiz", handlers.GetQuiz)
	r.Post("/api/quiz/check", handlers.CheckQuiz)
	return r
}
