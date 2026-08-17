package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"

	"github.com/NicolasFerreras/Gator/internal/errors_handling"
	"github.com/NicolasFerreras/Gator/internal/models"
)

func FetchFeed(ctx context.Context, feedURL string) (*models.RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, errors_handling.ErrMakeRequest(err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors_handling.ErrServer(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors_handling.ErrReadingResponse(err)
	}

	var feed models.RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, errors_handling.ErrUnmarshal(err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
