package repositories

import "github.com/wan6sta/go-url/internal/config"

type Repositories interface {
	CreateURL(URL string, baseURL string) (string, error)
	GetURL(ID string) (string, error)
}

type Repository struct {
	rs  Repositories
	cfg *config.Config
}

func NewRepository(rs Repositories, cfg *config.Config) *Repository {
	return &Repository{rs: rs, cfg: cfg}
}

func (r *Repository) CreateURL(URL string) (string, error) {
	return r.rs.CreateURL(URL, r.cfg.BaseURL)
}

func (r *Repository) GetURL(ID string) (string, error) {
	return r.rs.GetURL(ID)
}
