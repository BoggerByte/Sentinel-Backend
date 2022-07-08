package memdb

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

type Store interface {
	GetOauth2Flow(ctx context.Context, state string) (Oauth2Flow, error)
	SetOauth2Flow(ctx context.Context, state string, oauth2Flow Oauth2Flow, duration time.Duration) error
	DeleteOauth2Flow(ctx context.Context, state string) error
}

type Redis struct {
	client *redis.Client
}

type ConnectionConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisStore(cfg ConnectionConfig) (Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &Redis{
		client: client,
	}, nil
}
