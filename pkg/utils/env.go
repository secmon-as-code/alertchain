package utils

import (
	"os"
	"strings"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

func Env() types.EnvVars {
	vars := types.EnvVars{}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		vars[types.EnvVarName(pair[0])] = types.EnvVarValue(pair[1])
	}
	return vars
}
