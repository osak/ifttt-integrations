package mastodon

import (
	"context"
	"github.com/mattn/go-mastodon"
	"os"
)

type Client struct {
	client *mastodon.Client
}

func NewMastodonClient() *Client {
	clientKey := os.Getenv("MASTODON_CLIENT_KEY")
	clientSecret := os.Getenv("MASTODON_CLIENT_SECRET")
	accessToken := os.Getenv("MASTODON_ACCESS_TOKEN")

	config := &mastodon.Config{
		Server:       "https://social.mikutter.hachune.net",
		ClientID:     clientKey,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	}
	return &Client{
		client: mastodon.NewClient(config),
	}
}

func (c *Client) Post(message string) error {
	toot := &mastodon.Toot{
		Status: message,
	}
	_, err := c.client.PostStatus(context.Background(), toot)
	return err
}
