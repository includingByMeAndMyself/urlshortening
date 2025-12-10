package main

import (
	"log"
	"net/http"

	"github.com/includingByMeAndMyself/urlshortening/internal/config"
	"github.com/includingByMeAndMyself/urlshortening/internal/handler"
	"github.com/includingByMeAndMyself/urlshortening/internal/repository"
)

func main() {
	cfg := config.Load()

	repo, err := repository.NewFileRepo(cfg.FileStoragePath)
	if err != nil {
		log.Fatalf("failed to create file repository: %v", err)
	}
	log.Printf("Using file storage: %s", cfg.FileStoragePath)

	s := handler.NewServerWithRepo(repo, cfg.BaseURL)
	handler := s.Router()

	log.Printf("Starting server on %s, base URL: %s", cfg.Address, cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.Address, handler))
}
