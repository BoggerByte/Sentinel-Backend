package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrExpiredToken = errors.New("token has expired")

type Payload struct {
	ID            uuid.UUID `json:"id"`
	UserDiscordID string    `json:"user_discord_id"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiredAt     time.Time `json:"expired_at"`
}

func NewPayload(userDiscordID string, duration time.Duration) *Payload {
	id, _ := uuid.NewRandom()
	return &Payload{
		ID:            id,
		UserDiscordID: userDiscordID,
		IssuedAt:      time.Now(),
		ExpiredAt:     time.Now().Add(duration),
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
