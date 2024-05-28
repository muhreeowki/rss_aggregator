package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// urlToFeed takes a url and returns an RSSFeed struct
func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	// Get the RSS feed
	res, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	// Read the response body
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	feed := RSSFeed{}

	// Unmarshal the data into the feed struct
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return RSSFeed{}, err
	}
	return feed, nil
}
