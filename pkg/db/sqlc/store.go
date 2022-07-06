package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

type ConnectionConfig struct {
	Driver   string
	Protocol string
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func Init(config ConnectionConfig) (*SQLStore, error) {
	conn, err := sql.Open(config.Driver, fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		config.Protocol, config.Username, config.Password, config.Host, config.Port, config.Name, config.SSLMode))
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return &SQLStore{
		db:      conn,
		Queries: New(conn),
	}, nil
}

func (s *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
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
		return txErr
	}

	return tx.Commit()
}
