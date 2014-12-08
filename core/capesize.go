package core

type CloudProvider interface {
	InstanceConfig() ProviderMachineConfig
	CreateInstance(*Machine, chan bool)
}

type Capesize struct {
	Provider CloudProvider
}

func (c *Capesize) SetProvider(provider CloudProvider) {
	c.Provider = provider
}
