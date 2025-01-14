package utils

import (
	"encoding/json"

	"github.com/m-mizutani/goerr/v2"
)

func ToAny(input any) (any, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to marshal data", goerr.V("input", input))
	}

	var output any
	if err := json.Unmarshal(raw, &output); err != nil {
		return nil, goerr.Wrap(err, "Fail to unmarshal data", goerr.V("raw", string(raw)))
	}

	return output, nil
}
