package main

import (
	"fmt"
	"log"

	"github.com/phpires/gator/internal/config"
)

func main() {
	username := "MrPiresz"
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error while trying reading config: %v", err)
	}

	err = cfg.SetUser(username)
	if err != nil {
		log.Fatalf("error while trying to set user config: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error while trying reading config: %v", err)
	}

	fmt.Printf("db_url: %v\nusername: %v\n", cfg.DbUrl, cfg.CurrentUserName)
}
