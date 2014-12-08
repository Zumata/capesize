package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/Zumata/capesize/config"
	"github.com/Zumata/capesize/core"
	"github.com/Zumata/capesize/services"
	flags "github.com/jessevdk/go-flags"
)

type Options struct{}

var options Options
var parser = flags.NewParser(&options, flags.Default)

// Command: create

type CreateCommand struct {
}

var createCommand CreateCommand

func (c *CreateCommand) Execute(args []string) error {
	if len(args) != 2 {
		return errors.New("Please provide the provider and number of hosts to be spawned")
	}

	numHosts, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("Please provide the number of hosts to be spawned")
	}

	provider := args[0]
	err = config.FindProvider(provider)
	if err != nil {
		return errors.New("Specified provider is unknown.")
	}

	capesize := &core.Capesize{}

	success := services.CreateInstances(capesize, provider, numHosts)
	if !success {
		return errors.New("Deployment Failed")
	}
	return nil
}

// Command: list

type ListCommand struct {
}

var listCommand ListCommand

func (c *ListCommand) Execute(args []string) error {

	if len(args) != 1 {
		return errors.New("Please provide the provider")
	}

	provider := args[0]
	err := config.FindProvider(provider)
	if err != nil {
		return errors.New("Specified provider is unknown.")
	}

	capesize := &core.Capesize{}

	services.ListInstances(capesize, provider)

	return nil
}

func main() {

	parser.AddCommand("create",
		"Spawn docker hosts",
		"The create command spawns a docker host. capesize create <provider> <num hosts>.",
		&createCommand)

	parser.AddCommand("list",
		"List spawned docker hosts",
		"The list command lists spawned docker hosts. capesize list.",
		&listCommand)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

}
