package config

import "errors"

var RegisteredProviders = make(map[string]bool)

func RegisterProvider(provider string, validity bool) {
	RegisteredProviders[provider] = validity
}

func FindProvider(provider string) error {
	config, found := RegisteredProviders[provider]
	if !found {
		return errors.New("capesize: " + provider + " is not supported")
	}
	if !config {
		return errors.New("capesize: " + provider + " has not been appropriately configured")
	}
	return nil
}
