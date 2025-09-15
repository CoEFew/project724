package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"my-app-backend/internal/db"
)

type Feedback struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Message string `json:"message"`
	Source  string `json:"source"` // optional: ติด tag หน้า/ฟีเจอร์
}

// POST /api/feedback
func SaveFeedback(w http.ResponseWriter, r *http.Request) {
	var f Feedback
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	f.Name = strings.TrimSpace(f.Name)
	f.Contact = strings.TrimSpace(f.Contact)
	f.Message = strings.TrimSpace(f.Message)
	f.Source = strings.TrimSpace(f.Source)

	if len(f.Message) < 1 {
		http.Error(w, "message required", http.StatusBadRequest)
		return
	}
	if len(f.Message) > 4000 {
		f.Message = f.Message[:4000]
	}
	if len(f.Name) > 64 {
		f.Name = f.Name[:64]
	}
	if len(f.Contact) > 128 {
		f.Contact = f.Contact[:128]
	}
	if len(f.Source) > 64 {
		f.Source = f.Source[:64]
	}

	if err := db.InsertFeedback(r.Context(), f.Name, f.Contact, f.Message, f.Source); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

// GET /api/feedback?limit=5
func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	limit := 5
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	rows, err := db.GetRecentFeedbacks(r.Context(), limit)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type outRow struct {
		Name      string `json:"name"`
		Contact   string `json:"contact"`
		Message   string `json:"message"`
		Source    string `json:"source"`
		CreatedAt string `json:"created_at"`
	}

	out := make([]outRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, outRow{
			Name:      row.Name,
			Contact:   row.Contact,
			Message:   row.Message,
			Source:    row.Source,
			CreatedAt: row.CreatedAt.Format(timeLayout), // ใช้ layout ร่วม
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// timeLayout ใช้ร่วม (ถ้ามีไฟล์อื่นประกาศอยู่แล้วข้ามได้)
const timeLayout = "2006-01-02T15:04:05Z07:00"
