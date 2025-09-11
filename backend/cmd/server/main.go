package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"my-app-backend/internal/db"
	httpSrv "my-app-backend/internal/http"
)

func loadEnv() {
	candidates := []string{
		".env.local",
		"../.env.local",
		"../../.env.local",
		".env",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			log.Println("loaded env from", p)
			return
		}
	}
}

func main() {
	loadEnv()

	ctx := context.Background()
	if err := db.Init(ctx); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer db.Close() // << ใช้อันนี้เท่านั้น ไม่เรียก pool โดยตรง

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := httpSrv.NewRouter()
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
