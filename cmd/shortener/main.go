package main

import (
	"log"
	"net/http"

	"github.com/includingByMeAndMyself/urlshortening/internal/config"
	"github.com/includingByMeAndMyself/urlshortening/internal/handler"
	"github.com/includingByMeAndMyself/urlshortening/internal/repository"
	"github.com/includingByMeAndMyself/urlshortening/pkg/middleware"
)

func main() {
	cfg := config.MustLoad()

	repo, err := repository.NewFileRepo(cfg.FileStoragePath)
	if err != nil {
		log.Fatalf("failed to create file repository: %v", err)
	}
	log.Printf("Using file storage: %s", cfg.FileStoragePath)

	s := handler.NewServerWithRepo(repo, cfg.BaseURL)
	router := s.Router()

	handler := middleware.GzipDecompress(router)
	handler = middleware.Logging(handler)
	handler = middleware.GzipCompress(handler)

	log.Printf("Starting server on %s, base URL: %s", cfg.Address, cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.Address, handler))
}
