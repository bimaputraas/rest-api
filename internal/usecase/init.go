package usecase

import (
	"github.com/bimaputraas/rest-api/internal/config"
	"github.com/bimaputraas/rest-api/internal/repository"
)

type (
	Usecase struct {
		repo   repository.Repository
		config *config.Config
	}
)

func New(repo repository.Repository, config *config.Config) *Usecase {
	return &Usecase{
		repo:   repo,
		config: config,
	}
}
