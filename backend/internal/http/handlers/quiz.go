package handlers

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"my-app-backend/internal/db"
)

// ==== (1) Data & Config ====

type QuizResp struct {
	ID        string `json:"id"`
	HintCount int    `json:"hintCount"`
	Token     string `json:"token"` // HMAC(secret, answer|id|exp)
	Exp       int64  `json:"exp"`   // unix seconds
}

type HintReq struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
	Index int    `json:"index"` // 1 หรือ 2
}

// ดึง secret จาก env
func getSecret() []byte {
	sec := os.Getenv("HMAC_SECRET")
	if sec == "" {
		sec = "change-me-in-env"
	}
	return []byte(sec)
}

// สุ่ม id ปลอดภัย
func randomID() (string, error) {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// สร้าง HMAC(base64) จาก s
func sign(secret []byte, s string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(s))
	sum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

// เทียบ HMAC
func equalHMAC(a, b string) bool {
	ab, _ := base64.StdEncoding.DecodeString(a)
	bb, _ := base64.StdEncoding.DecodeString(b)
	return hmac.Equal(ab, bb)
}

// ==== (2) Helpers ====

func poolTierForLevel(level int) int {
	// ใช้ชุดยากตั้งแต่ level >= 4
	if level >= 4 && level < 10 {
		return 2
	}
	if level >= 10 {
		return 3
	}
	return 1
}

// ==== (3) Handlers ====

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	level := 1
	if lv := r.URL.Query().Get("level"); lv != "" {
		if n, err := strconv.Atoi(lv); err == nil && n > 0 {
			level = n
		}
	}
	tier := poolTierForLevel(level)

	// Get category from query parameter
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "สัตว์" // Default to animals
	}

	now := time.Now()
	rand.Seed(now.UnixNano())

	var q db.QuizRow
	var err error

	// Try to get quiz by category first, fallback to any tier if not found
	if category != "" {
		q, err = db.GetRandomQuizByTierAndCategory(r.Context(), tier, category)
		if err != nil && errors.Is(err, db.ErrNoQuiz) {
			// Fallback to any category if specific category not found
			q, err = db.GetRandomQuizByTier(r.Context(), tier)
		}
	} else {
		q, err = db.GetRandomQuizByTier(r.Context(), tier)
	}

	if err != nil {
		if errors.Is(err, db.ErrNoQuiz) {
			http.Error(w, "no_quiz", http.StatusNotFound)
			return
		}
		log.Printf("GetQuiz: query error: %v", err)
		http.Error(w, "cannot get quiz", http.StatusInternalServerError)
		return
	}

	id, err := randomID()
	if err != nil {
		http.Error(w, "cannot generate id", http.StatusInternalServerError)
		return
	}
	exp := now.Add(60 * time.Second).Unix()

	secret := getSecret()
	payload := q.Answer + "|" + id + "|" + strconv.FormatInt(exp, 10)
	token := sign(secret, payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store") // ← กัน cache

	resp := QuizResp{
		ID:        id,
		Token:     token,
		Exp:       exp,
		HintCount: 2,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

type CheckReq struct {
	ID    string `json:"id"`
	Guess string `json:"guess"`
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

func CheckQuiz(w http.ResponseWriter, r *http.Request) {
	var req CheckReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// หมดอายุ?
	if time.Now().Unix() > req.Exp {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"correct": false, "reason": "expired"})
		return
	}

	// เทียบ HMAC(secret, guess|id|exp) กับ token
	secret := getSecret()
	normalizedGuess := normalizeThai(req.Guess)
	signGuess := sign(secret, normalizedGuess+"|"+req.ID+"|"+strconv.FormatInt(req.Exp, 10))
	ok := equalHMAC(signGuess, req.Token)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"correct": ok})
}

// helper: normalize
func normalizeThai(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

// ==== เฉลยหลังหมดเวลา ====

type RevealReq struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

func RevealQuiz(w http.ResponseWriter, r *http.Request) {
	var req RevealReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// เฉลยได้เมื่อหมดเวลาแล้วเท่านั้น
	if time.Now().Unix() < req.Exp {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"answer": "", "error": "not_expired"})
		return
	}

	secret := getSecret()

	answers, err := db.GetAllAnswers(r.Context())
	if err != nil {
		http.Error(w, "cannot reveal", http.StatusInternalServerError)
		return
	}

	var found string
	for _, a := range answers {
		if equalHMAC(sign(secret, a+"|"+req.ID+"|"+strconv.FormatInt(req.Exp, 10)), req.Token) {
			found = a
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"answer": found})
}

func GetHint(w http.ResponseWriter, r *http.Request) {
	var req HintReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Index != 1 && req.Index != 2 {
		http.Error(w, "invalid index", http.StatusBadRequest)
		return
	}
	// หมดเวลาแล้วห้ามขอเพิ่ม
	if time.Now().Unix() > req.Exp {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "expired"})
		return
	}

	secret := getSecret()

	// ดึง answer+hint1+hint2 ทั้งหมด (หรือทำฟังก์ชันหาเฉพาะที่ match ก็ได้)
	rows, err := db.GetAllQuizLite(r.Context()) // ← ดูโค้ดส่วน db ด้านล่าง
	if err != nil {
		http.Error(w, "cannot get hint", http.StatusInternalServerError)
		return
	}

	var hint string
	for _, row := range rows {
		signed := sign(secret, row.Answer+"|"+req.ID+"|"+strconv.FormatInt(req.Exp, 10))
		if equalHMAC(signed, req.Token) {
			if req.Index == 1 {
				hint = row.Hint1
			} else {
				hint = row.Hint2
			}
			break
		}
	}

	if hint == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"index": req.Index,
		"hint":  hint,
	})
}

// ==== Multiple Choice Quiz Endpoint ====

type MultipleChoiceQuizResp struct {
	ID        string `json:"id"`
	Answer    string `json:"answer"`
	Hint1     string `json:"hint1"`
	Hint2     string `json:"hint2"`
	HintCount int    `json:"hintCount"`
	Token     string `json:"token"`
	Exp       int64  `json:"exp"`
}

// GetQuizForMultipleChoice returns quiz data with answer for multiple-choice games
func GetQuizForMultipleChoice(w http.ResponseWriter, r *http.Request) {
	levelStr := r.URL.Query().Get("level")
	category := r.URL.Query().Get("category")

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 {
		level = 1
	}

	ctx := r.Context()
	quizzes, err := db.GetQuizzes(ctx, level, category)
	if err != nil {
		log.Printf("Error fetching quizzes: %v", err)
		http.Error(w, "cannot fetch quiz", http.StatusInternalServerError)
		return
	}

	if len(quizzes) == 0 {
		log.Printf("No quizzes found for level=%d, category=%s", level, category)
		http.Error(w, "no quiz available", http.StatusNotFound)
		return
	}

	// Randomly select a quiz
	quiz := quizzes[rand.Intn(len(quizzes))]

	// Generate secure ID and token
	id, err := randomID()
	if err != nil {
		log.Printf("Error generating ID: %v", err)
		http.Error(w, "cannot generate quiz", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	exp := now.Add(60 * time.Second).Unix()

	// Create token for this quiz
	secret := getSecret()
	tokenData := quiz.Answer + "|" + id + "|" + strconv.FormatInt(exp, 10)
	token := sign(secret, tokenData)

	// Return quiz data with answer for multiple-choice games
	resp := MultipleChoiceQuizResp{
		ID:        id,
		Answer:    quiz.Answer,
		Hint1:     quiz.Hint1,
		Hint2:     quiz.Hint2,
		HintCount: 2,
		Token:     token,
		Exp:       exp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
