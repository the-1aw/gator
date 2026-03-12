package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/the-1aw/gator/internal/database"
)

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
