package memdb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *Redis) SetSession(ctx context.Context, session Session, duration time.Duration) (Session, error) {
	key := fmt.Sprintf("session_%s", session.ID)
	return session, r.client.Set(ctx, key, &session, duration).Err()
}

func (r *Redis) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	key := fmt.Sprintf("session_%s", id)
	c := r.client.Get(ctx, key)
	if err := c.Err(); err != nil {
		return Session{}, err
	}

	var session Session
	if err := c.Scan(&session); err != nil {
		return Session{}, err
	}
	return session, nil
}
