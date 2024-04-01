package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/utils"
	"log/slog"
	"os"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLNotValid = errors.New("url not valid")
)

type Storage struct {
	cfg *config.Config
	log *slog.Logger
}

type URLItem struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func NewStorage(cfg *config.Config, log *slog.Logger) *Storage {
	return &Storage{cfg: cfg, log: log}
}

func (s *Storage) CreateURL(URL string, baseURL string) (string, error) {
	if URL == "" {
		s.log.Error("url not valid", ErrURLNotValid)
		return "", ErrURLNotValid
	}

	ShortUrl := utils.GenerateID(6)

	Uuid, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("cannot generate uuid", err)
		return "", err
	}

	finalURL := baseURL + "/" + ShortUrl

	url := URLItem{
		Uuid:        Uuid.String(),
		ShortUrl:    ShortUrl,
		OriginalUrl: URL,
	}

	data, err := json.Marshal(url)
	if err != nil {
		s.log.Error("cannot marshal url", err)
		return "", err
	}

	_, err = os.Stat(s.cfg.StoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(s.cfg.StoragePath)
			if err != nil {
				s.log.Error("error create file", err)
				return "", err
			}
			defer file.Close()
			slog.Info("created storage file")

			_, err = file.Write(data)
			if err != nil {
				s.log.Error("write to file", err)
				return "", err
			}
		} else {
			s.log.Error("error stat file", err)
			return "", err
		}
	} else {
		file, err := os.OpenFile(s.cfg.StoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			s.log.Error("error open file", err)
			return "", err
		}
		defer file.Close()

		_, err = file.WriteString("\n")
		if err != nil {
			s.log.Error("error write to file", err)
			return "", err
		}

		_, err = file.Write(data)
		if err != nil {
			s.log.Error("error write to file", err)
			return "", err
		}
	}

	return finalURL, nil
}

func (s *Storage) GetURL(ID string) (string, error) {
	file, err := os.Open(s.cfg.StoragePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return "", err
	}
	defer file.Close()

	var URLItems []URLItem
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var URL URLItem
		if err := json.Unmarshal([]byte(scanner.Text()), &URL); err != nil {
			fmt.Println("Ошибка при распаковке JSON:", err)
			continue
		}

		URLItems = append(URLItems, URL)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return "", err
	}

	for _, u := range URLItems {
		if u.ShortUrl == ID {
			return u.OriginalUrl, nil
		}
	}

	return "", ErrURLNotFound
}
