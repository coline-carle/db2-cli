package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var db2Cmd = &cobra.Command{
	Use:   "db2-cli",
	Short: "World of Warcraft db2 toolkit",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func addCommands() {
	db2Cmd.AddCommand(headerCmd)
}

// Execute adds all child commands to the root command
func Execute() {
	addCommands()

	if c, err := db2Cmd.ExecuteC(); err != nil {
		c.Println("")
		c.Println(c.UsageString())
		os.Exit(-1)
	}
}
