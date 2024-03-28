package storage

import (
	"errors"
	"github.com/wan6sta/go-url/internal/app/utils"
)

var (
	ErrUrlNotFound = errors.New("url not found")
)

type Storage struct {
	urls map[string]string
}

func NewStorage() *Storage {
	return &Storage{urls: make(map[string]string)}
}

func (s *Storage) CreateUrl(url string, baseUrl string) (string, error) {
	ID := utils.GenerateID(6)

	s.urls[ID] = url

	return baseUrl + "/" + ID, nil
}

func (s *Storage) GetUrl(ID string) (string, error) {
	url, ok := s.urls[ID]

	if !ok {
		return "", ErrUrlNotFound
	}

	return url, nil
}
