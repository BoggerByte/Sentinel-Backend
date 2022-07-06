package util

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func GenerateSnowflakeID() snowflake.ID {
	node, _ := snowflake.NewNode(1)
	return node.Generate()
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[RandomInt(0, len(charset)-1)]
	}
	return string(b)
}
