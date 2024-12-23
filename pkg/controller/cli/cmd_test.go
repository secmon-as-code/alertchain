package cli_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/controller/cli"
)

func TestCLI(t *testing.T) {
	ctx := context.Background()
	gt.NoError(t, cli.New().Run(ctx, []string{"alertchain"}))
}
