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

	s := handler.NewServer()
	router := s.Router()

	loggedRouter := middleware.Logging(router)

	log.Printf("Starting server on %s", cfg.Address)
	log.Fatal(http.ListenAndServe(cfg.Address, loggedRouter))
}
