package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mattn/go-mastodon"
	"log"
)

type Payload struct {
	Comment string
	Place   string
	Url     string
}

func shareMastodon(message string) error {
	client := NewMastodonClient()
	toot := &mastodon.Toot{
		Status: message,
	}
	_, err := client.PostStatus(context.Background(), toot)
	return err
}

func HandleRequest(payload Payload) error {
	log.Printf("%v", payload)

	message := fmt.Sprintf("%s (@ %s) %s", payload.Comment, payload.Place, payload.Url)
	if err := shareMastodon(message); err != nil {
		log.Printf("Failed to post to Mastodon: %v", err)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
