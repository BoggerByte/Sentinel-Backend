package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (m *PasetoMaker) CreateToken(userDiscordID string, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(userDiscordID, duration)
	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	return token, payload, err
}

func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload = new(Payload)
	err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
