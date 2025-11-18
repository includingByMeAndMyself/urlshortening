package service

import (
	"crypto/rand"
	"errors"
	"strings"

	"github.com/includingByMeAndMyself/urlshortening/internal/model"
	"github.com/includingByMeAndMyself/urlshortening/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Shorten(originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", errors.New("invalid URL")
	}

	id, err := generateShortID()
	if err != nil {
		return "", err
	}

	pair := &model.URLPair{
		ID:          id,
		OriginalURL: originalURL,
	}
	s.repo.Save(pair)
	return id, nil
}

func (s *Service) GetOriginal(id string) (string, error) {
	pair, ok := s.repo.Get(id)
	if !ok {
		return "", errors.New("not found")
	}
	return pair.OriginalURL, nil
}

func isValidURL(u string) bool {
	u = strings.TrimSpace(u)
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

func generateShortID() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i, v := range b {
		b[i] = charset[v%byte(len(charset))]
	}
	return string(b), nil
}
