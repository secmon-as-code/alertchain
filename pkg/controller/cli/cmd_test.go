package cli_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/controller/cli"
	"github.com/m-mizutani/gt"
)

func TestCLI(t *testing.T) {
	ctx := context.Background()
	gt.NoError(t, cli.New().Run(ctx, []string{"alertchain"}))
}
