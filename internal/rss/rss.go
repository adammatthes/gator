package rss

import (
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"html"
)

type RSSItem struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	myRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(myRequest)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := &RSSFeed{}
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}

	unescapeRoutine(result)

	return result, nil
}

func unescapeRoutine(r *RSSFeed) {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	for i, _ := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(r.Channel.Item[i].Title)
		r.Channel.Item[i].Description = html.UnescapeString(r.Channel.Item[i].Description)
	}
}
