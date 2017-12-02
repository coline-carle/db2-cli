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
	stringsCmd.RunE = strings
}

func strings(cmd *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}

	wdb6, err := wdb6.Decode(f)
	if err != nil {
		return err
	}

	if !wdb6.Header.HasStringTable() {
		return fmt.Errorf("the file has no string table")
	}

	return nil
}
