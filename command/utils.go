package command

import (
	"os"

	"github.com/urfave/cli"
)

func fileParam(c *cli.Context) (*os.File, error) {
	if c.NArg() != 1 {
		return nil, cli.NewExitError("invalid number of arguments", 1)
	}
	filename := c.Args().Get(0)
	return os.Open(filename)
}
