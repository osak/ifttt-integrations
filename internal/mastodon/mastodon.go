package mastodon

import (
	"context"
	"os"

	"github.com/mattn/go-mastodon"
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

func (c *Client) PostImages(message string, filenames []string) error {
	mediaIDs := make([]mastodon.ID, len(filenames))
	for i := 0; i < len(filenames); i += 1 {
		id, err := c.UploadImage(filenames[i])
		if err != nil {
			return err
		}
		mediaIDs[i] = id
	}

	toot := &mastodon.Toot{
		Status:   message,
		MediaIDs: mediaIDs,
	}
	_, err := c.client.PostStatus(context.Background(), toot)
	return err
}

func (c *Client) UploadImage(filename string) (mastodon.ID, error) {
	attachment, err := c.client.UploadMedia(context.Background(), filename)
	return attachment.ID, err
}
