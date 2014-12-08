package main

import (
	"os"
	"strconv"

	"github.com/Zumata/capesize/config"
	"github.com/Zumata/capesize/core"
	"github.com/Zumata/capesize/services"
)

func main() {

	if len(os.Args) != 3 {
		panic("Please provide the provider and number of hosts to be spawned")
	}

	numHosts, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Please provide the number of hosts to be spawned")
	}

	provider := os.Args[1]
	err = config.FindProvider(provider)
	if err != nil {
		panic(err)
	}

	capesize := &core.Capesize{}
	success := services.CreateInstances(capesize, provider, numHosts)
	if !success {
		panic("Deployment Failed")
	}

}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
