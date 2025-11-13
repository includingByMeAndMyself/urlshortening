package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/includingByMeAndMyself/urlshortening/internal/model"
	"github.com/includingByMeAndMyself/urlshortening/internal/repository"
)

type mockRepo struct {
	urls map[string]*model.URLPair
}

func newMockRepo() repository.Repository {
	return &mockRepo{
		urls: make(map[string]*model.URLPair),
	}
}

func (m *mockRepo) Save(pair *model.URLPair) {
	m.urls[pair.ID] = pair
}

func (m *mockRepo) Get(id string) (*model.URLPair, bool) {
	pair, ok := m.urls[id]
	return pair, ok
}

const testBaseURL = "http://localhost:8080"

func TestShortenHandler(t *testing.T) {
	repo := newMockRepo()
	svc := NewServerWithRepo(repo, testBaseURL)
	router := svc.Router()

	tests := []struct {
		name        string
		body        string
		contentType string
		wantStatus  int
		wantPrefix  string
	}{
		{
			name:        "valid URL",
			body:        "https://practicum.yandex.ru/",
			contentType: "text/plain",
			wantStatus:  http.StatusCreated,
			wantPrefix:  testBaseURL + "/",
		},
		{
			name:        "invalid URL",
			body:        "not-a-url",
			contentType: "text/plain",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "empty body",
			body:        "",
			contentType: "text/plain",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "missing Content-Type",
			body:        "https://example.com",
			contentType: "",
			wantStatus:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantPrefix != "" && !strings.HasPrefix(w.Body.String(), tt.wantPrefix) {
				t.Errorf("expected body to start with %q, got %q", tt.wantPrefix, w.Body.String())
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	repo := newMockRepo()

	testID := "test123"
	testURL := "https://practicum.yandex.ru/"
	repo.Save(&model.URLPair{
		ID:          testID,
		OriginalURL: testURL,
	})

	svc := NewServerWithRepo(repo, testBaseURL)
	router := svc.Router()

	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "valid redirect",
			path:         "/" + testID,
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: testURL,
		},
		{
			name:       "non-existent ID",
			path:       "/unknown",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "root GET (invalid)",
			path:       "/",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantLocation != "" {
				loc := w.Header().Get("Location")
				if loc != tt.wantLocation {
					t.Errorf("expected Location %q, got %q", tt.wantLocation, loc)
				}
			}
		})
	}
}
