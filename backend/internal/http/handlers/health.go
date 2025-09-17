// internal/http/handlers/health.go
package handlers

import (
	"encoding/json"
	"net/http"

	"my-app-backend/internal/db"
)

type healthResp struct {
	Status string `json:"status"`
	DB     string `json:"db,omitempty"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// พยายาม ping DB อย่างสุภาพ (timeout ควรไปตั้งที่ context upstream ถ้าต้องการ)
	if err := db.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthResp{
			Status: "down",
			DB:     "down",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResp{
		Status: "ok",
		DB:     "ok",
	})
}
