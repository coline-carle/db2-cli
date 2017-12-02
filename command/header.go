package command

import (
	"github.com/wow-sweetlie/db2-cli/formats/wdb6"

	"github.com/urfave/cli"
)

// Header display db2 header
func Header(c *cli.Context) error {
	f, err := fileParam(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	header, err := wdb6.Decode(f)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	header.Print()
	return nil
}
