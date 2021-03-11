package main

import (
	"fmt"
	"ifttt-integrations/internal/mastodon"
	twWrapper "ifttt-integrations/internal/twitter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"gopkg.in/errgo.v2/fmt/errors"
)

type Payload struct {
	TweetUrl string
}

func downloadImage(url string) (string, error) {
	extPos := strings.LastIndex(url, ".")
	if extPos == -1 {
		return "", errors.Newf("Cannot determine image extension from url: %s", url)
	}
	ext := url[extPos+1:]

	pattern := fmt.Sprintf("switch-crosspost-*.%s", ext)
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", err
	}

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	io.Copy(f, res.Body)
	f.Close()

	return f.Name(), nil
}

func trimUrls(text string) string {
	re := regexp.MustCompile("http\\S+")
	return re.ReplaceAllString(text, "")
}

func HandleRequest(payload Payload) error {
	log.Printf("%v", payload)

	re := regexp.MustCompile(".*status/(\\d+)$")
	matches := re.FindStringSubmatch(payload.TweetUrl)
	if matches == nil {
		log.Printf("Failed to extract tweet id from URL: %s", payload.TweetUrl)
		return nil
	}

	tweetId, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		log.Printf("Failed to parse tweet id as a number: %d", tweetId)
		return nil
	}

	tw := twWrapper.NewTwitterClient()
	params := twitter.StatusShowParams{}
	tweet, _, err := tw.Statuses.Show(tweetId, &params)
	if err != nil {
		log.Printf("Failed to get tweet %d", tweetId)
		return nil
	}

	images := make([]string, len(tweet.ExtendedEntities.Media))
	for i := 0; i < len(tweet.ExtendedEntities.Media); i += 1 {
		url := tweet.ExtendedEntities.Media[i].MediaURLHttps
		filename, err := downloadImage(url)
		if err != nil {
			log.Printf("Failed to download image: %s", url)
			return nil
		}
		images[i] = filename
	}

	message := trimUrls(tweet.Text)

	md := mastodon.NewMastodonClient()
	if err := md.PostImages(message, images); err != nil {
		log.Printf("Failed to post to Mastodon: %v", err)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
