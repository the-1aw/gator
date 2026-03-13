package cli

import (
	"context"
	"errors"
	"fmt"
)

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
	user, err := s.db.CreateUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Println(user)
	return nil
}

func handleReset(s *state, _ command) error {
	return s.db.DeleteAllUsers(context.Background())
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
