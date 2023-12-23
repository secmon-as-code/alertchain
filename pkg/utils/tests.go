package utils

import (
	"os"

	"github.com/m-mizutani/goerr"
)

type EnvLoader func() error

func Env(key string, dst *string) EnvLoader {
	return func() error {
		v, ok := os.LookupEnv(key)
		if !ok {
			return goerr.New("No such env: %s", key)
		}
		*dst = v
		return nil
	}
}

func LoadEnv(envs ...EnvLoader) error {
	for _, env := range envs {
		if err := env(); err != nil {
			return err
		}
	}
	return nil
}
