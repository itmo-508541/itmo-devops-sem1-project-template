//go:build dist

package config

import "os"

func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
