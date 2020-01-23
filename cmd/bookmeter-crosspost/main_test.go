package main

import (
	"github.com/antchfx/htmlquery"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseReviewPage(t *testing.T) {
	f, err := os.Open("../../testdata/bookmeter_review.html")
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(f)
	if err != nil {
		panic(err)
	}

	info, err := parseReviewPage(doc)
	if err != nil {
		t.Fatalf("Parse must success: %v", err)
	}

	assert.Equal(t, "いま世界の哲学者が考えていること", info.title)
	assert.Equal(t, "岡本 裕一朗", info.author)
	assert.Equal(t, "現代の哲学における主要なトピックを広く浅くまとめた本。読み口はかなり軽いので、とりあえずざっと読んで興味の糸口を探すといいと思う。「ポストモダン以後」の思想をまとめた第1章は自分が最近考えていることに合っていて、次に掘り下げるべきキーワードが分かったので良かった。", info.review)
}

func TestParseBookPage(t *testing.T) {
	f, err := os.Open("../../testdata/bookmeter_book.html")
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(f)
	if err != nil {
		panic(err)
	}

	info, err := parseBookPage(doc)
	if err != nil {
		t.Fatalf("Parse must success: %v", err)
	}

	assert.Equal(t, "異種族レビュアーズ　4 (ドラゴンコミックスエイジ)", info.title)
	assert.Equal(t, "masha", info.author)
	assert.Equal(t, "を読んだ", info.review)
}
