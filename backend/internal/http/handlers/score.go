package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
)

type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

var (
	scores []Score
	mu     sync.Mutex
)

// POST /api/scores
func SaveScore(w http.ResponseWriter, r *http.Request) {
	var s Score
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mu.Lock()
	scores = append(scores, s)
	sort.Slice(scores, func(i, j int) bool { return scores[i].Score > scores[j].Score })
	if len(scores) > 10 {
		scores = scores[:10]
	}
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

// GET /api/scores
func GetScores(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}
