package command

import (
	"fmt"
	"os"

	"github.com/wow-sweetlie/db2-cli/formats/wdb6"

	"github.com/urfave/cli"
)

// Header display db2 header
func Header(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("invalid number of arguments", 33)
	}
	filename := c.Args().Get(0)
	f, err := os.Open(filename)
	if err != nil {
		return cli.NewExitError("error opening file", 33)
	}

	header, err := wdb6.Decode(f)
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError("error deccoding file", 33)
	}
	wdb6.PrintHeader(header)
	return nil
}
