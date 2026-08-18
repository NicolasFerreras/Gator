package rss

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Tech &amp; News</title>
    <link>http://example.com</link>
    <description>Un feed de prueba</description>
    <item>
      <title>Post &amp; More</title>
      <link>http://example.com/p1</link>
      <description>Desc 1</description>
      <pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
    </item>
    <item>
      <title>Segundo post</title>
      <link>http://example.com/p2</link>
      <description>Desc 2</description>
      <pubDate>Tue, 03 Jan 2006 15:04:05 MST</pubDate>
    </item>
  </channel>
</rss>`

func TestFetchFeed_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, sampleRSS)
	}))
	defer srv.Close()

	feed, err := FetchFeed(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("FetchFeed: %v", err)
	}
	if feed.Channel.Title != "Tech & News" {
		t.Fatalf("Channel.Title = %q, want %q (html unescape)", feed.Channel.Title, "Tech & News")
	}
	if len(feed.Channel.Item) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(feed.Channel.Item))
	}
	if feed.Channel.Item[0].Title != "Post & More" {
		t.Fatalf("Item[0].Title = %q, want %q", feed.Channel.Item[0].Title, "Post & More")
	}
}

func TestFetchFeed_InvalidXML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "esto no es xml válido <<<")
	}))
	defer srv.Close()

	_, err := FetchFeed(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("se esperaba error de unmarshal con XML inválido")
	}
}

func TestFetchFeed_ServerUnreachable(t *testing.T) {
	// Puerto cerrado: fuerza error de conexión sin red externa.
	_, err := FetchFeed(context.Background(), "http://127.0.0.1:1/feed.xml")
	if err == nil {
		t.Fatal("se esperaba error de conexión a servidor caído")
	}
}
