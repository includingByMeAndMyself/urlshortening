package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/includingByMeAndMyself/urlshortening/internal/repository"
	"github.com/includingByMeAndMyself/urlshortening/internal/service"
)

type Server struct {
	service *service.Service
	baseURL string
}

func NewServer(baseURL string) *Server {
	repo := repository.NewMemoryRepo()
	return NewServerWithRepo(repo, baseURL)
}

func NewServerWithRepo(repo repository.URLStorer, baseURL string) *Server {
	svc := service.New(repo)
	return &Server{
		service: svc,
		baseURL: baseURL,
	}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/", s.handleShorten)
	r.Get("/{id}", s.handleRedirect)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "root GET not allowed", http.StatusBadRequest)
	})
	return r
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}

	id, err := s.service.Shorten(originalURL)
	if err != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	shortURL := s.baseURL + "/" + id
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(shortURL)); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	originalURL, err := s.service.GetOriginal(id)
	if err != nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
