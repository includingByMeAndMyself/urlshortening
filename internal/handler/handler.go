package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/includingByMeAndMyself/urlshortening/internal/repository"
	"github.com/includingByMeAndMyself/urlshortening/internal/service"
)

type Server struct {
	service *service.Service
}

func NewServer() *Server {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)
	return &Server{service: svc}
}

func NewServerWithRepo(repo repository.Repository) *Server {
	svc := service.New(repo)
	return &Server{service: svc}
}

func (s *Server) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleMain)
	return mux
}

func (s *Server) handleMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.handleShorten(w, r)
	} else if r.Method == http.MethodGet {
		path := r.URL.Path
		if path == "/" {
			http.Error(w, "root GET not allowed", http.StatusBadRequest)
			return
		}
		s.handleRedirect(w, r)
	} else {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
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

	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	originalURL, err := s.service.GetOriginal(id)
	if err != nil {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect) // 307
}
