package config

import (
	"fmt"
)

func OptionalEnv(key string, defaultValue string) string {
	value, ok := lookupEnv(key)
	if !ok {
		value = defaultValue
	}

	return value
}

func RequiredEnv(key string) string {
	value, ok := lookupEnv(key)
	if !ok {
		panic(fmt.Errorf("env.%s is required", key))
	}

	return value
}
