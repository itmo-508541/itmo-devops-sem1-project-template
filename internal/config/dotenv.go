//go:build !dist

package config

import (
	"os"

	"github.com/joho/godotenv"
)

func lookupEnv(key string) (string, bool) {
	for _, dotEnv := range []string{".env", ".env.local"} {
		if _, err := os.Stat(dotEnv); err == nil {
			if err := godotenv.Overload(dotEnv); err != nil {
				panic(err)
			}
		}
	}

	return os.LookupEnv(key)
}
