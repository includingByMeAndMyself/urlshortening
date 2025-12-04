package repository

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/includingByMeAndMyself/urlshortening/internal/model"
)

type urlRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileRepo struct {
	filePath string
	urls     map[string]*model.URLPair
	mu       sync.RWMutex
}

func NewFileRepo(filePath string) (*FileRepo, error) {
	repo := &FileRepo{
		filePath: filePath,
		urls:     make(map[string]*model.URLPair),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *FileRepo) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var records []urlRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}

	for _, record := range records {
		r.urls[record.ShortURL] = &model.URLPair{
			ID:          record.ShortURL,
			OriginalURL: record.OriginalURL,
		}
	}

	return nil
}

func (r *FileRepo) save() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := make([]urlRecord, 0, len(r.urls))
	for id, pair := range r.urls {
		records = append(records, urlRecord{
			UUID:        id,
			ShortURL:    id,
			OriginalURL: pair.OriginalURL,
		})
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := r.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, r.filePath)
}

func (r *FileRepo) Save(pair *model.URLPair) {
	r.mu.Lock()
	r.urls[pair.ID] = pair
	r.mu.Unlock()

	if err := r.save(); err != nil {
	}
}

func (r *FileRepo) Get(id string) (*model.URLPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pair, ok := r.urls[id]
	return pair, ok
}
