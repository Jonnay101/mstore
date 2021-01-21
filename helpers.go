package mstore

import (
	"fmt"
	"os"
)

// MustEnv errors if the environment variable cannot be found
func MustEnv(key string) (val string, err error) {
	val = os.Getenv(key)
	if val == "" {
		err = fmt.Errorf("environment variable '%s', not found", key)
	}
	return
}
