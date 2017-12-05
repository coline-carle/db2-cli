package commands

import (
	"os"

	"github.com/wow-sweetlie/db2-cli/formats/wdb6"

	"github.com/spf13/cobra"
)

var (
	showFields bool
)

// HeaderCmd for main
var headerCmd = &cobra.Command{
	Use:   "header [db2 file(s)]",
	Short: "Print db2 file header",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	headerCmd.Run = header
}

// Header display db2 header
func header(cmd *cobra.Command, args []string) {
	f, err := os.Open(args[0])
	checkErr(err)

	header, err := wdb6.DecodeHeader(f)
	checkErr(err)

	wdb6.PrintHeader(header)
}
