package main

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/net/html"
	"ifttt-integrations/internal/mastodon"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Payload struct {
	Url string
}

type bookInfo struct {
	title  string
	author string
	url    string
	review string
}

func mustClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("Failed to close: %v", err)
	}
}

func parseReviewPage(doc *html.Node) (bookInfo, error) {
	data := htmlquery.FindOne(doc, "//div[@class=\"bm-page-loader\"]/@data-resource")
	if data == nil {
		return bookInfo{}, fmt.Errorf("cannot find data blob")
	}
	rawJson, err := url.QueryUnescape(htmlquery.InnerText(data))
	if err != nil {
		return bookInfo{}, fmt.Errorf("cannot decode blob: %w", err)
	}

	var blob struct {
		Path     string `json:"path"`
		Content  string `json:"content"`
		Contents struct {
			Book struct {
				Title  string `json:"title"`
				Author struct {
					Name string `json:"name"`
				} `json:"author"`
			} `json:"book"`
		} `json:"contents"`
	}
	dec := json.NewDecoder(strings.NewReader(rawJson))
	if err := dec.Decode(&blob); err != nil {
		return bookInfo{}, fmt.Errorf("cannot decode json: %w", err)
	}

	return bookInfo{
		title:  blob.Contents.Book.Title,
		author: blob.Contents.Book.Author.Name,
		url:    fmt.Sprintf("https://bookmeter.com%s", blob.Path),
		review: blob.Content,
	}, nil
}

func parseBookPage(doc *html.Node) (bookInfo, error) {
	nodes := htmlquery.Find(doc, "//ul[@class=\"breadcrumb-list\"]/li")
	if len(nodes) < 2 {
		return bookInfo{}, fmt.Errorf("breadcrumb list has less than 2 elements")
	}
	titleNode := nodes[len(nodes)-1]
	authorNode := nodes[len(nodes)-2]

	urlNode := htmlquery.FindOne(doc, "//meta[@property=\"og:url\"]/@content")
	if urlNode == nil {
		return bookInfo{}, fmt.Errorf("cannot find og:url")
	}

	return bookInfo{
		title:  htmlquery.InnerText(titleNode),
		author: htmlquery.InnerText(authorNode),
		url:    htmlquery.InnerText(urlNode),
		review: "を読んだ",
	}, nil
}

func HandleRequest(payload Payload) error {
	log.Printf("%v", payload)

	res, err := http.Get(payload.Url)
	if err != nil {
		return fmt.Errorf("failed to load bookmeter page: %w", err)
	}
	defer mustClose(res.Body)

	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse bookmeter page: %w", err)
	}

	var info bookInfo
	if strings.Contains(payload.Url, "reviews") {
		info, err = parseReviewPage(doc)
	} else {
		info, err = parseBookPage(doc)
	}
	if err != nil {
		return fmt.Errorf("failed to extract info from bookmeter page: %w", err)
	}

	log.Printf("%v", info)

	message := fmt.Sprintf("【%s/%s】%s %s", info.title, info.author, info.review, payload.Url)
	client := mastodon.NewMastodonClient()
	if err := client.Post(message); err != nil {
		log.Printf("Failed to post to Mastodon: %v", err)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
