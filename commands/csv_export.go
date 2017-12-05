package commands

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/wow-sweetlie/db2-cli/formats/dblayout"
	"github.com/wow-sweetlie/db2-cli/formats/wdb6"

	"github.com/spf13/cobra"
)

var csvExportCmd = &cobra.Command{
	Use:   "csvexport [db2 file]",
	Short: "csvexport export the content of a db2 file to csv",
	Args:  cobra.MinimumNArgs(1),
}

var build int

func init() {
	csvExportCmd.Run = csvExport
	csvExportCmd.Flags().IntVarP(&build, "build", "b", 25549, "build version")
}

func csvExport(cmd *cobra.Command, args []string) {
	layout, err := dblayout.LoadLayout(build)
	checkErr(err)

	f, err := os.Open(args[0])
	checkErr(err)

	filename := path.Base(args[0])
	tablename := strings.Split(filename, ".")[0]

	for _, table := range layout.Tables {
		if table.Name == tablename && build == table.Build {
			line := make([]string, len(table.Fields))

			for i, field := range table.Fields {
				line[i] = field.Name
			}

			records, err := wdb6.Decode(f, table)
			checkErr(err)

			w := csv.NewWriter(os.Stdout)
			w.Write(line)
			checkErr(err)

			for _, record := range records {
				for i, field := range record {
					line[i] = fmt.Sprintf("%v", field)
				}
				w.Write(line)
				checkErr(err)
			}
			w.Flush()
		}
	}

}
