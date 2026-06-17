package main

import (
	"log"
	"os"

	"github.com/phpires/gator/internal/config"
)

func main() {
	var appState state
	var appCommands commands

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error while trying reading config: %v", err)
	}

	appState.configState = &cfg
	appCommands.handlers = map[string]func(*state, command) error{}
	appCommands.handlers["login"] = handlerLogin

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatalf("Takes at least one argument to execute gator")
	}

	userCmd := command{
		name: userArgs[1],
		args: userArgs[2:],
	}

	err = appCommands.run(&appState, userCmd)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}
