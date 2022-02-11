package cmd

import "github.com/urfave/cli/v2"

func Run(argv []string) error {
	cli := cli.App{}

	if err := cli.Run(argv); err != nil {
		return err
	}
	return nil
}
