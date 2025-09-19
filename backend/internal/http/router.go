package http

import (
	"net/http"
	"os"
	"time"

	"my-app-backend/internal/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// ---------- Base middlewares ----------
	r.Use(middleware.RequestID)                 // ใส่ X-Request-ID ให้ตามรอยง่าย
	r.Use(middleware.RealIP)                    // ดึง IP จริงหลัง CDN/Proxy
	r.Use(middleware.Recoverer)                 // กันแอปล้มจาก panic
	r.Use(middleware.Timeout(15 * time.Second)) // กันแฮงค์ (รวมทั้ง preflight/options)

	// ---------- CORS ----------
	// อนุญาต FE หลัก + localhost ระหว่าง dev
	fe := os.Getenv("ALLOWED_ORIGIN")
	if fe == "" {
		fe = "https://project724-frontend.onrender.com"
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{fe, "http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// ใช้ * เพื่อเลี่ยงปัญหา header แปลก ๆ จากเบราว์เซอร์/axios ใน preflight
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false, // ไม่ใช้คุกกี้ข้ามโดเมน
		MaxAge:           300,
	}))

	// ---------- Health ----------
	// Render จะเรียกเช็กเป็นระยะ และคุณก็ใช้ดูสถานะได้
	r.Get("/health", handlers.Health)
	// หรือจะตอบ 204 ไวมาก ๆ ก็ทำเพิ่มได้:
	r.Head("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })

	// ---------- Quiz / Scores / Chat / Feedback ----------
	r.Get("/api/quiz", handlers.GetQuiz)
	r.Post("/api/quiz/reveal", handlers.RevealQuiz)
	r.Post("/api/quiz/check", handlers.CheckQuiz)
	r.Post("/api/quiz/hint", handlers.GetHint)

	r.Post("/api/scores", handlers.SaveScore)
	r.Get("/api/scores", handlers.GetScores)

	r.Post("/api/chat", handlers.ChatHandler)
	r.Post("/api/feedback", handlers.SaveFeedback)
	r.Get("/api/feedback", handlers.GetFeedbacks)

	// ---------- Party mode ----------
	r.Route("/api/rooms", func(r chi.Router) {
		r.Post("/", handlers.CreateRoom)          // POST /api/rooms
		r.Post("/{code}/join", handlers.JoinRoom) // POST /api/rooms/{code}/join
		r.Post("/{code}/ready", handlers.ReadyRoom)
		r.Post("/{code}/start", handlers.StartRoom)
		r.Post("/{code}/guess", handlers.GuessRoom)
		r.Post("/{code}/leave", handlers.LeaveRoom)
		r.Get("/{code}/snapshot", handlers.RoomSnapshot)
	})

	// WebSocket (CORS ไม่บังคับใช้กับ WS; ตัว upgrader.CheckOrigin(true) อยู่ใน handlers แล้ว)
	r.Get("/ws/rooms/{code}", handlers.RoomWS)

	return r
}
