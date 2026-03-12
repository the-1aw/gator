package cli

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
