package commands

import (
	"fmt"
	"os"

	"github.com/wow-sweetlie/db2-cli/formats/wdb6"

	"github.com/spf13/cobra"
)

var stringsCmd = &cobra.Command{
	Use:   "strings [db2 file(s)]",
	Short: "extract strings of a db2 file",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	stringsCmd.Run = exportStrings
}

func exportStrings(cmd *cobra.Command, args []string) {
	f, err := os.Open(args[0])
	checkErr(err)

	positions, strings, err := wdb6.DecodeStrings(f)
	checkErr(err)

	for _, position := range positions {
		fmt.Printf("%08x: ", position)
		fmt.Println(strings[position])
	}
}
