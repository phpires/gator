package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/phpires/gator/internal/config"
	"github.com/phpires/gator/internal/database"
)

type state struct {
	configState *config.Config
	dbState     *database.Queries
}

func main() {
	var appState state
	var appCommands commands

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error while trying reading config: %v", err)
	}

	appState.configState = &cfg
	appCommands.handlers = map[string]func(*state, command) error{}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerListUsers)
	appCommands.register("agg", handlerFetch)
	appCommands.register("addfeed", middlewareLoggedIn(handlerAddFeeds))
	appCommands.register("feeds", handlerListFeeds)
	appCommands.register("follow", middlewareLoggedIn(handlerFeedFollow))
	appCommands.register("following", middlewareLoggedIn(handlerFollowing))
	appCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error while trying to connect to dabase: %v", err)
	}
	appState.dbState = database.New(db)

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatalf("Takes at least one argument to execute gator")
	}

	userCmd := command{
		Name: userArgs[1],
		Args: userArgs[2:],
	}

	err = appCommands.run(&appState, userCmd)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}
