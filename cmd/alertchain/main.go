package main

import (
	"os"

	"github.com/m-mizutani/alertchain/pkg/controller"
)

func main() {
	controller.New().CLI(os.Args)
}
