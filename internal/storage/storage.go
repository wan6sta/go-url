package storage

import (
	"errors"
	"github.com/wan6sta/go-url/internal/utils"
)

var (
	ErrURLNotFound = errors.New("url not found")
)

type Storage struct {
	urls map[string]string
}

func NewStorage() *Storage {
	return &Storage{urls: make(map[string]string)}
}

func (s *Storage) CreateURL(URL string, baseUrl string) (string, error) {
	ID := utils.GenerateID(6)

	s.urls[ID] = URL

	return baseUrl + "/" + ID, nil
}

func (s *Storage) GetURL(ID string) (string, error) {
	URL, ok := s.urls[ID]

	if !ok {
		return "", ErrURLNotFound
	}

	return URL, nil
}
