package utils

import (
	"os"

	"github.com/m-mizutani/goerr/v2"
)

type EnvLoader func() error

func EnvDef(key string, dst *string) EnvLoader {
	return func() error {
		v, ok := os.LookupEnv(key)
		if !ok {
			return goerr.New("No such env: " + key)
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
