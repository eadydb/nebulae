package envs

import (
	"os"

	"github.com/eadydb/nebulae/pkg/utils/format"
)

// EnvPrefix is the key prefix of the environment variable value
var EnvPrefix = "KWOK_"

// GetEnv returns the value of the environment variable named by the key.
func GetEnv[T any](key string, def T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	if value == "" {
		return def
	}
	t, err := format.Parse[T](value)
	if err != nil {
		return def
	}
	return t
}

// GetEnvWithPrefix returns the value of the environment variable named by the key with kwok prefix.
func GetEnvWithPrefix[T any](key string, def T) T {
	return GetEnv(EnvPrefix+key, def)
}
