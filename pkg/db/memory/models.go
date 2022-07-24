package memdb

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	DiscordID    string    `json:"discord_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

type Oauth2Flow struct {
	Completed     bool  `json:"completed"`
	UserDiscordID int64 `json:"user_discord_id"`
}

func (f *Oauth2Flow) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Oauth2Flow) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}
