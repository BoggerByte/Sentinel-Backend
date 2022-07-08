package token

import (
	"github.com/BoggerByte/Sentinel-backend.git/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPasetoMaker(t *testing.T) {
	testCases := []struct {
		name         string
		symmetricKey string
		checkResult  func(t *testing.T, maker Maker, err error)
	}{
		{
			name:         "OK",
			symmetricKey: util.RandomString(32),
			checkResult: func(t *testing.T, maker Maker, err error) {
				require.NotEmpty(t, maker)
				require.NoError(t, err)
			},
		},
		{
			name:         "WrongKeyLength",
			symmetricKey: "definitely_not_32_characters",
			checkResult: func(t *testing.T, maker Maker, err error) {
				require.Empty(t, maker)
				require.Error(t, err)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewPasetoMaker(tc.symmetricKey)
			tc.checkResult(t, maker, err)
		})
	}
}

func TestPasetoTokenMaker(t *testing.T) {
	userDiscordID := util.RandomSnowflakeID().Int64()
	issuedAt := time.Now()

	testCases := []struct {
		name          string
		scope         string
		userDiscordID int64
		duration      time.Duration
		checkVerify   func(t *testing.T, payload *Payload, err error)
	}{
		{
			name:          "OK",
			userDiscordID: userDiscordID,
			duration:      time.Minute,
			checkVerify: func(t *testing.T, payload *Payload, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, payload)
				require.NotZero(t, payload.ID)
				require.Equal(t, userDiscordID, payload.UserDiscordID)
				require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
				require.WithinDuration(t, issuedAt.Add(time.Minute), payload.ExpiredAt, time.Second)
			},
		},
		{
			name:          "TokenExpired",
			userDiscordID: userDiscordID,
			duration:      -time.Minute,
			checkVerify: func(t *testing.T, payload *Payload, err error) {
				require.ErrorIs(t, err, ErrExpiredToken)
				require.Empty(t, payload)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewPasetoMaker(util.RandomString(32))
			require.NoError(t, err)

			accessToken, _, err := maker.CreateToken(userDiscordID, tc.duration)
			require.NoError(t, err)
			require.NotEmpty(t, accessToken)

			payload, err := maker.VerifyToken(accessToken)
			tc.checkVerify(t, payload, err)
		})
	}
}
