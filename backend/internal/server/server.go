package server

import (
	"log"
	"my-app-backend/internal/config"
	internalhttp "my-app-backend/internal/http"
	"net/http"
)

func Run() {
	config.LoadEnv()
	router := internalhttp.NewRouter()
	log.Printf("Server running on port %s", config.GetPort())
	http.ListenAndServe(":"+config.GetPort(), router)
}
