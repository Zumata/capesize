package services

import (
	"fmt"

	"github.com/Zumata/capesize/amazon"
	"github.com/Zumata/capesize/core"
)

func CreateInstances(cs *core.Capesize, providerName string, numHosts int) (success bool) {

	if providerName == "amazon" {
		cs.SetProvider(amazon.NewAmazon(amazon.DefaultEC2Config))
	}

	machines := []*core.Machine{}
	for i := 0; i < numHosts; i++ {
		machines = append(machines, core.NewMachine(cs.Provider.InstanceConfig()))
	}

	// Launch machines and wait until all have been configured

	launchChan := make(chan bool)

	for _, m := range machines {
		go cs.Provider.CreateInstance(m, launchChan)
	}

	success = false
	for i := 0; i < len(machines); i++ {
		result := <-launchChan
		if result {
			success = true
		}
		fmt.Printf("Waiting for servers %d / %d\n", i+1, len(machines))
	}

	return
}
