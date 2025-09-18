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
	r.Post("/api/quiz/reveal", handlers.RevealQuiz)
	r.Post("/api/quiz/check", handlers.CheckQuiz)
	r.Post("/api/scores", handlers.SaveScore)
	r.Get("/api/scores", handlers.GetScores)
	r.Post("/api/chat", handlers.ChatHandler)
	r.Post("/api/quiz/hint", handlers.GetHint)
	r.Post("/api/feedback", handlers.SaveFeedback)
	r.Get("/api/feedback", handlers.GetFeedbacks)

	r.Route("/api/rooms", func(r chi.Router) {
		r.Post("/", handlers.CreateRoom)          // POST /api/rooms
		r.Post("/{code}/join", handlers.JoinRoom) // POST /api/rooms/{code}/join
		r.Post("/{code}/ready", handlers.ReadyRoom)
		r.Post("/{code}/start", handlers.StartRoom)
		r.Post("/{code}/guess", handlers.GuessRoom)
		r.Post("/{code}/leave", handlers.LeaveRoom)
	})

	// WebSocket
	r.Get("/ws/rooms/{code}", handlers.RoomWS)

	return r
}
