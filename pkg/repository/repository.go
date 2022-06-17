package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Guilds interface {
}

type Configs interface {
}

type Repository struct {
	Authorization
	Guilds
	Configs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
