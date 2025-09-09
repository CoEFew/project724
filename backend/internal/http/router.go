package http

import (
	"my-app-backend/internal/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://project724-frontend.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/health", handlers.Health)
	r.Get("/api/quiz", handlers.GetQuiz)
	r.Post("/api/quiz/check", handlers.CheckQuiz)
	r.Post("/api/scores", handlers.SaveScore)
	r.Get("/api/scores", handlers.GetScores)
	r.Post("/api/chat", handlers.ChatHandler)
	return r
}
