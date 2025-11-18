package cli

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/the-1aw/gator/internal/config"
	"github.com/the-1aw/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

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

type commandRegistry struct {
	cmds map[string]func(*state, command) error
}

func newCommandRegistry() commandRegistry {
	return commandRegistry{cmds: make(map[string]func(*state, command) error)}
}

func (c *commandRegistry) run(s *state, cmd command) error {
	if handler, ok := c.cmds[cmd.name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("command \"%s\" not found", cmd.name)
}

func (c *commandRegistry) register(name string, handler func(*state, command) error) error {
	c.cmds[name] = handler
	return nil
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("agg only expect a single <DURATION> argument")
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("login only expect a single <USERNAME> argument")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return err
	}

	err := s.cfg.SetUser(cmd.args[0])
	fmt.Printf("Username has been set to: %s\n", cmd.args[0])
	return err
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("register only expect a single <NAME> argument")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Println(user)
	return nil
}

func handleReset(s *state, _ command) error {
	err := s.db.DeleteAllUsers(context.Background())
	return err
}

func handleUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		userTxt := user.Name
		if user.Name == s.cfg.CurrentUsername {
			userTxt += " (current)"
		}
		fmt.Printf("* %s\n", userTxt)
	}
	return nil
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("%s missing arguments\nusage: addfeed feedname url\n", cmd.name)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	fmt.Println(feed)
	return nil
}

func handleFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeedSummary(context.Background())
	if err == nil {
		for _, feed := range feeds {
			fmt.Printf("%s\n%s\n", feed.Name, feed.CreatedBy.String)
		}
	}
	return err
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("%s missing arguments\nusage: follow url\n", cmd.name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err == nil {
		for _, item := range ff {
			fmt.Println(item.FeedName, item.Username)
		}
	}
	return err
}

func handleFollowing(s *state, cmd command, user database.User) error {
	following, err := s.db.GetUserFeedFollow(context.Background(), user.Name)
	if err != nil {
		return nil
	}
	for _, follow := range following {
		fmt.Println(follow)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("%s missing arguments\nusage: follow url\n", cmd.name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
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

func handleBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(parsedLimit)
	} else {
		limit = 2
	}
	posts, err := s.db.GetPostForUser(context.Background(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err == nil {
		for _, post := range posts {
			fmt.Println(post)
		}
	}
	return err
}

func Run() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("Unable to read gator config:\n%s", err)
	}

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	s := state{cfg: &c, db: dbQueries}
	cr := newCommandRegistry()

	cr.register("login", handlerLogin)
	cr.register("register", handlerRegister)
	cr.register("reset", handleReset)
	cr.register("users", handleUsers)
	cr.register("agg", handleAgg)
	cr.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cr.register("feeds", handleFeeds)
	cr.register("follow", middlewareLoggedIn(handleFollow))
	cr.register("following", middlewareLoggedIn(handleFollowing))
	cr.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cr.register("browse", middlewareLoggedIn(handleBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments provided")
	}
	cmdName := args[1]
	cmdArgs := args[2:]
	err = cr.run(&s, command{name: cmdName, args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
