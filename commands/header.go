package commands

import (
	"fmt"
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
	headerCmd.Flags().BoolVarP(&showFields, "fields", "f", false, "show fields format")
	headerCmd.RunE = header
}

// Header display db2 header
func header(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}

	db, err := wdb6.Decode(f)
	if err != nil {
		return err
	}
	wdb6.PrintHeader(db.Header)
	if showFields {
		fmt.Print("\n\n")
		wdb6.PrintFieldsFormat(db.FieldsFormat)
	}
	return nil
}
