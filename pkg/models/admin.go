package models

import (
	"time"
)

type Admin struct {
	ID        int64
	DiscordID int64
	Nickname  string
	CreatedAt time.Time
}
