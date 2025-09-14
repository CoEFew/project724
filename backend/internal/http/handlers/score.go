package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"my-app-backend/internal/db"
)

type Score struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	GameName string `json:"gamename"`
}

// POST /api/scores
func SaveScore(w http.ResponseWriter, r *http.Request) {
	var s Score
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.Name = strings.TrimSpace(s.Name)
	s.GameName = strings.TrimSpace(s.GameName)

	if s.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if s.GameName == "" {
		http.Error(w, "gamename required", http.StatusBadRequest)
		return
	}
	if s.Score < 0 {
		s.Score = 0
	}
	if s.Score > 1_000_000 {
		s.Score = 1_000_000
	}

	if err := db.InsertScore(r.Context(), s.Name, s.Score, s.GameName); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GET /api/scores?limit=10&gamename=PolaJigsaw
func GetScores(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	game := strings.TrimSpace(r.URL.Query().Get("gamename")) // "" = ทุกเกม

	rows, err := db.GetTopScores(r.Context(), game, limit)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type outScore struct {
		Name     string `json:"name"`
		Score    int    `json:"score"`
		GameName string `json:"gamename"`
	}
	out := make([]outScore, 0, len(rows))
	for _, row := range rows {
		out = append(out, outScore{Name: row.Name, Score: row.Score, GameName: row.GameName})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
