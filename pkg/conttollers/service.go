package conttollers

import "github.com/BoggerByte/Sentinel-backend.git/pkg/repository"

type Authorization interface {
}

type Guilds interface {
}

type Configs interface {
}

type Service struct {
	Authorization
	Guilds
	Configs
}

func NewService(repository *repository.Repository) *Service {
	return &Service{}
}
