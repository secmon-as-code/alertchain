package cli_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/controller/cli"
	"github.com/m-mizutani/gt"
)

func TestCLI(t *testing.T) {
	gt.NoError(t, cli.New().Run([]string{"alertchain"}))
}
