package main

import (
	"log"
	"net/http"

	configs "cmd/auth/config"
	"cmd/auth/internal/handler"
	"cmd/auth/internal/service"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	userService := service.NewUserService(cfg)
	loginHandler := handler.NewLoginHandler(userService)
	verifyHandler := handler.NewVerifyHandler(userService)

	http.HandleFunc("/login", loginHandler.Handle)
	http.HandleFunc("/verify", verifyHandler.Handle)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
