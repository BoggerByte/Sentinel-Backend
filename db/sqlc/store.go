package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

var db *Store

type ConnectionConfig struct {
	Driver   string
	Source   string
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func Init(config ConnectionConfig) error {
	conn, err := sql.Open(config.Driver, fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		config.Source, config.Username, config.Password, config.Host, config.Port, config.Name, config.SSLMode))
	if err != nil {
		return err
	}
	db = &Store{
		db:      conn,
		Queries: New(conn),
	}
	return nil
}

func GetStore() *Store {
	return db
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	txErr := fn(q)
	if txErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", txErr, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateAdminTxParams struct {
	Admin  CreateAdminParams
	Guilds []CreateGuildParams `json:"guilds"`
}

type CreateAdminTxResult struct {
	Admin  Admin
	Guilds []Guild `json:"guilds"`
}

// CreateAdminTx creates admin, creates guilds, creates relation admin <> guilds
func (s *Store) CreateAdminTx(ctx context.Context, arg CreateAdminTxParams) (CreateAdminTxResult, error) {
	var result CreateAdminTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Admin, err = q.CreateAdmin(ctx, arg.Admin)
		if err != nil {
			return err
		}

		for _, guild := range arg.Guilds {
			rGuild, err := q.CreateGuild(ctx, guild)
			result.Guilds = append(result.Guilds, rGuild)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
