package repositories

import "github.com/wan6sta/go-url/internal/app/config"

type Repositories interface {
	CreateUrl(url string) (string, error)
	GetUrl(ID string) (string, error)
}

type Repository struct {
	rs  Repositories
	cfg *config.Config
}

func NewRepository(rs Repositories, cfg *config.Config) *Repository {
	return &Repository{rs: rs, cfg: cfg}
}

func (r *Repository) CreateUrl(url string) (string, error) {
	return r.rs.CreateUrl(url)
}

func (r *Repository) GetUrl(ID string) (string, error) {
	return r.rs.GetUrl(ID)
}
