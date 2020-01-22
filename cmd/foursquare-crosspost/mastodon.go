package main

import (
	"github.com/mattn/go-mastodon"
	"os"
)

func NewMastodonClient() *mastodon.Client {
	clientKey := os.Getenv("MASTODON_CLIENT_KEY")
	clientSecret := os.Getenv("MASTODON_CLIENT_SECRET")
	accessToken := os.Getenv("MASTODON_ACCESS_TOKEN")

	config := &mastodon.Config{
		Server:       "https://social.mikutter.hachune.net",
		ClientID:     clientKey,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	}
	return mastodon.NewClient(config)
}
