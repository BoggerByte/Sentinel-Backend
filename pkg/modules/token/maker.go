package token

import "time"

type Maker interface {
	CreateToken(userDiscordID string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
