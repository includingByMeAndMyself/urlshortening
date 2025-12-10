package repository

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/includingByMeAndMyself/urlshortening/internal/model"
)

type urlRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// FileRepo представляет файловый репозиторий как обёртку над MemoryRepo
type FileRepo struct {
	*MemoryRepo
	filePath string
	mu       sync.Mutex
}

func NewFileRepo(filePath string) (*FileRepo, error) {
	repo := &FileRepo{
		MemoryRepo: NewMemoryRepo(),
		filePath:   filePath,
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *FileRepo) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		// Если файл не существует, это нормально - создаём пустой репозиторий
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var records []urlRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}

	// Загружаем записи в память через встроенный MemoryRepo
	for _, record := range records {
		r.MemoryRepo.Save(&model.URLPair{
			ID:          record.ShortURL,
			OriginalURL: record.OriginalURL,
		})
	}

	return nil
}

func (r *FileRepo) save() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Получаем все записи из встроенного MemoryRepo
	urls := r.MemoryRepo.GetAll()
	records := make([]urlRecord, 0, len(urls))
	for id, pair := range urls {
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
	// Сначала сохраняем в память через встроенный MemoryRepo
	r.MemoryRepo.Save(pair)

	// Затем сохраняем в файл (защищено мьютексом)
	if err := r.save(); err != nil {
		log.Printf("failed to save to file: %v", err)
	}
}
