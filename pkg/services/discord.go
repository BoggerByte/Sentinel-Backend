package services

import (
	"context"
	"encoding/json"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"net/url"
)

type DiscordOauth2Service struct {
	config *oauth2.Config
}

func NewDiscordOauth2Service(config *oauth2.Config) *DiscordOauth2Service {
	return &DiscordOauth2Service{config: config}
}

func (s *DiscordOauth2Service) NewURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *DiscordOauth2Service) NewInviteBotURL() string {
	v := url.Values{}
	v.Set("client_id", s.config.ClientID)
	v.Set("permissions", "8") // administrator
	v.Set("scope", discord.ScopeBot)
	return s.config.Endpoint.AuthURL + "?" + v.Encode()
}

func (s *DiscordOauth2Service) Exchange(code string) (*oauth2.Token, error) {
	return s.config.Exchange(context.Background(), code)
}

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int64  `json:"flags"`
	Banner        string `json:"banner"`
	AccentColor   int64  `json:"accent_color"`
	PublicFlags   int64  `json:"public_flags"`
}

func (s *DiscordOauth2Service) GetUser(token *oauth2.Token) (DiscordUser, error) {
	resp, err := s.config.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil {
		return DiscordUser{}, err
	}

	var discordUser DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		return DiscordUser{}, err
	}
	return discordUser, nil
}

type DiscordGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	IsOwner     bool   `json:"owner"`
	Permissions int64  `json:"permissions"`
}

func (s *DiscordOauth2Service) GetUserGuilds(token *oauth2.Token) ([]DiscordGuild, error) {
	resp, err := s.config.Client(context.Background(), token).Get("https://discord.com/api/users/@me/guilds")
	if err != nil {
		return []DiscordGuild{}, err
	}

	var discordGuilds []DiscordGuild
	if err := json.NewDecoder(resp.Body).Decode(&discordGuilds); err != nil {
		return []DiscordGuild{}, err
	}
	return discordGuilds, nil
}
