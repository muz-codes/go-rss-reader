package go_rss_reader

import (
	"encoding/xml"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *RssItem) Parse(rssUrls []string) ([]RssItem, error) {
	var allFeedItems []RssItem
	ch := make(chan interface{})
	for _, singleUrl := range rssUrls {
		go feedReader(singleUrl, ch)
	}
	for i := 0; i < len(rssUrls); i++ {
		output := <-ch
		switch output.(type) {
		case error:
			logger.Error("error in Parse", zap.Error(output.(error)))
		case []RssItem:
			allFeedItems = append(allFeedItems, output.([]RssItem)...)
		}
	}
	return allFeedItems, nil
}

func feedReader(url string, ch chan interface{}) {
	resp, err := callUrl(url)
	if err != nil {
		logger.Error("error in Parse", zap.Error(err))
		ch <- err
		return
	}
	defer resp.Body.Close()
	rssXmlVar := rssXml{}

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rssXmlVar)
	if err != nil {
		logger.Error(fmt.Sprintf("error in Parse for %v", url), zap.Error(err))
		ch <- err
		return
	}
	var items = make([]RssItem, len(rssXmlVar.Channel.Items))
	for index, channelItem := range rssXmlVar.Channel.Items {
		assignValuesToRssItem(&items[index], rssXmlVar.Channel, channelItem)
	}
	ch <- items
	return
}

func assignValuesToRssItem(item *RssItem, channel feedChannel, channelItem item) {
	item.Source = channel.Source
	item.SourceURL = channel.SourceUrl
	item.Title = channelItem.Title
	item.Link = channelItem.Link
	item.Description = channelItem.Desc
	item.PublishDate = refinePubDate(channelItem.PubDate)
}

func callUrl(address string) (*http.Response, error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("error in callUrl while making request for %v", address), zap.Error(err))
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("error in callUrl while calling the request for %v", address), zap.Error(err))
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorOccurred := errors.New(fmt.Sprintf("response code is %v", resp.StatusCode))
		logger.Error(fmt.Sprintf("error in callUrl having non-success response for %v", address), zap.Error(errorOccurred))
		return nil, err
	}
	return resp, nil
}

func refinePubDate(pubDate string) time.Time {
	refinedDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", pubDate)
	if err != nil {
		refinedDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", pubDate)
		if err != nil {
			logger.Error(fmt.Sprintf("error in refinePubDate for %v", pubDate), zap.Error(err))
		}
		return refinedDate
	}
	return refinedDate
}
