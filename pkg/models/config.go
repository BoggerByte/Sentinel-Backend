package models

import "time"

type Config struct {
	ID        int64
	Version   int64
	Filename  string
	CreatedAt time.Time

	GuildID int64
}
