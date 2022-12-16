package go_rss_reader

import "time"

type RssItem struct {
	Title       string    `json:"title"`
	Source      string    `json:"source"`
	SourceURL   string    `json:"source_url"`
	Link        string    `json:"link"`
	PublishDate time.Time `json:"publish_date"`
	Description string    `json:"description"`
}

type item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Desc    string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

type feedChannel struct {
	Source    string `xml:"title"`
	SourceUrl string `xml:"link"`
	Items     []item `xml:"item"`
}

type rssXml struct {
	Channel feedChannel `xml:"channel"`
}
