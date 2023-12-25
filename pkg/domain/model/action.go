package model

import (
	"encoding/json"

	"github.com/m-mizutani/goerr"
)

type ActionArgs map[string]any

func ArgDef[T any](key string, dst *T, options ...ArgOption) ArgParser {
	var opt argParserOption
	for _, o := range options {
		o(&opt)
	}

	return func(args ActionArgs) error {
		v, ok := args[key]
		if !ok {
			if opt.Optional {
				return nil
			}
			return goerr.New("No such Optional key in action args").With("key", key)
		}

		raw, err := json.Marshal(v)
		if err != nil {
			return goerr.Wrap(err, "Failed to marshal action args").With("key", key)
		}

		var src T
		if err := json.Unmarshal(raw, &src); err != nil {
			return goerr.Wrap(err, "Failed to unmarshal action args").With("key", key)
		}

		*dst = src

		return nil
	}
}

type argParserOption struct {
	Optional bool
}

type ArgOption func(*argParserOption)

func ArgOptional() ArgOption {
	return func(opt *argParserOption) {
		opt.Optional = true
	}
}

type ArgParser func(args ActionArgs) error

func (x ActionArgs) Parse(psr ...ArgParser) error {
	for _, p := range psr {
		if err := p(x); err != nil {
			return err
		}
	}
	return nil
}
