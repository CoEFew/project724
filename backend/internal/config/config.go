package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	godotenv.Load()
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
