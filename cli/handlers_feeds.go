package cli

import (
	"context"
	"fmt"

	"github.com/the-1aw/gator/internal/database"
)

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("%s missing arguments\nusage: addfeed feedname url\n", cmd.name)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   cmd.args[0],
		Url:    cmd.args[1],
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
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
		UserID: user.ID,
		FeedID: feed.ID,
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
