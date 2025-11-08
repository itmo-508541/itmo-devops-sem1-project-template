package config

import (
	"fmt"
)

func optionalEnv(key string, defaultValue string) string {
	value, ok := lookupEnv(key)
	if !ok {
		value = defaultValue
	}

	return value
}

func requiredEnv(key string) string {
	value, ok := lookupEnv(key)
	if !ok {
		panic(fmt.Errorf("env.%s is required", key))
	}

	return value
}
