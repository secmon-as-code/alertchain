package main

import (
	"context"
	"os"

	"github.com/secmon-lab/alertchain/pkg/controller/cli"
)

func main() {
	ctx := context.Background()
	if err := cli.New().Run(ctx, os.Args); err != nil {
		os.Exit(1)
	}
}
