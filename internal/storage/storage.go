package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/xlillium/go-postman-tui/internal/domain"
)

type Storage struct {
	mutex    sync.Mutex
	filePath string
	requests []domain.Request
}

func NewStorage(filePath string) *Storage {
	storage := &Storage{
		filePath: filePath,
	}
	_ = storage.Load() 
	return storage
}

func (s *Storage) Load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		s.requests = []domain.Request{}
		return nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.requests)
}

func (s *Storage) Save() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(s.requests, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *Storage) AddRequest(req domain.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.requests = append(s.requests, req)
}

func (s *Storage) GetRequests() []domain.Request {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.requests
}
