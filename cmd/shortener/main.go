package main

import (
	"log"
	"net/http"

	"github.com/includingByMeAndMyself/urlshortening/internal/config"
	"github.com/includingByMeAndMyself/urlshortening/internal/handler"
	"github.com/includingByMeAndMyself/urlshortening/pkg/middleware"
)

func main() {
	cfg := config.MustLoad()

	s := handler.NewServer(cfg.BaseURL)
	router := s.Router()

	handler := middleware.GzipDecompress(router)
	handler = middleware.Logging(handler)
	handler = middleware.GzipCompress(handler)

	log.Printf("Starting server on %s, base URL: %s", cfg.Address, cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.Address, handler))
}
