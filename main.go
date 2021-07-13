package main

import (
	"os"

	"github.com/m-mizutani/alertchain/pkg/controller/cli"
)

func main() {
	ctrl := cli.New()
	if err := ctrl.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
