package main

import (
	"os"

	"github.com/wow-sweetlie/db2-cli/command"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "header",
			Usage:  "display db2 header",
			Action: command.Header,
		},
	}

	app.Run(os.Args)
}
