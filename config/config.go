package config

import "os"

func SetConfig(envName string, defaultVal string) string {
	envVal := os.Getenv(envName)
	if envVal == "" {
		return defaultVal
	}
	return envVal
}
