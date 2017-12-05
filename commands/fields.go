package commands

import (
	"os"

	"github.com/wow-sweetlie/db2-cli/formats/wdb6"

	"github.com/spf13/cobra"
)

var fieldsCmd = &cobra.Command{
	Use:   "fields [db2 file(s)]",
	Short: "Print db2 fields format",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	fieldsCmd.Run = fields
}

// Header display db2 header
func fields(cmd *cobra.Command, args []string) {
	f, err := os.Open(args[0])
	checkErr(err)

	fieldsFormat, err := wdb6.DecodeFieldsFormat(f)
	checkErr(err)

	wdb6.PrintFieldsFormat(fieldsFormat)
}
