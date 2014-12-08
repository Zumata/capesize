package services

import (
	"github.com/Zumata/capesize/amazon"
	"github.com/Zumata/capesize/core"
)

func ListInstances(cs *core.Capesize, providerName string) {

	if providerName == "amazon" {
		cs.SetProvider(amazon.NewAmazon(amazon.DefaultEC2Config))
	}

	cs.Provider.ListInstances()

}
