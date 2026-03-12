package cli

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/the-1aw/gator/internal/database"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func (r *RSSFeed) unescape() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	for idx := range r.Channel.Items {
		r.Channel.Items[idx].Title = html.UnescapeString(r.Channel.Items[idx].Title)
		r.Channel.Items[idx].Description = html.UnescapeString(r.Channel.Items[idx].Description)
	}
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	feed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &feed, err
	}
	req.Header.Set("user-agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &feed, err
	}
	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return &feed, err
	}
	err = xml.Unmarshal(data, &feed)
	if err == nil {
		feed.unescape()
	}
	return &feed, err
}

func scrapeFeed(s *state) {
	storedFeed, _ := s.db.GetNextFeedToFetch(context.Background())
	s.db.MarkFeedFetched(context.Background(), storedFeed.ID)
	fetchedFeed, _ := fetchFeed(context.Background(), storedFeed.Url)
	for _, item := range fetchedFeed.Channel.Items {
		pubDate, timeParseErr := time.Parse(time.RFC822, item.PubDate)
		_, postCreationErr := s.db.CreatePost(context.Background(), database.CreatePostParams{
			FeedID:      storedFeed.ID,
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: timeParseErr == nil},
			Url:         item.Link,
		})
		if postCreationErr != nil {
			log.Fatal(postCreationErr)
		}
	}
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("agg only expect a single <DURATION> argument")
	}
	fetchInterval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(fetchInterval)
	fmt.Printf("Collecting feeds every %s\n", fetchInterval.String())
	for ; ; <-ticker.C {
		scrapeFeed(s)
	}
}
