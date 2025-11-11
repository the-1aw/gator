package main

import (
	"fmt"
	"log"

	"github.com/the-1aw/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("Unable to read gator config:\n%s", err)
	}
	c.SetUser("lane")
	c, err = config.Read()
	if err != nil {
		log.Fatalf("Unable to read gator config:\n%s", err)
	}
	fmt.Println(c)
}
