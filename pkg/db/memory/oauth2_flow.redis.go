package memdb

import (
	"context"
	"fmt"
	"time"
)

func (r *Redis) SetOauth2Flow(ctx context.Context, state string, oauth2Flow Oauth2Flow, duration time.Duration) error {
	key := fmt.Sprintf("state_%s", state)
	return r.client.Set(ctx, key, &oauth2Flow, duration).Err()
}

func (r *Redis) GetOauth2Flow(ctx context.Context, state string) (Oauth2Flow, error) {
	key := fmt.Sprintf("state_%s", state)
	c := r.client.Get(ctx, key)
	if err := c.Err(); err != nil {
		return Oauth2Flow{}, err
	}

	var oauth2flow Oauth2Flow
	if err := c.Scan(&oauth2flow); err != nil {
		return Oauth2Flow{}, err
	}
	return oauth2flow, nil
}

func (r *Redis) DeleteOauth2Flow(ctx context.Context, state string) error {
	key := fmt.Sprintf("state_%s", state)
	return r.client.Del(ctx, key).Err()
}
