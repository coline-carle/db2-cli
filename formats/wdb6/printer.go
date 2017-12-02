package wdb6

import (
	"fmt"
	"github.com/fatih/color"
)

var (
	fieldNameColor = color.New(color.FgCyan).Add(color.Bold)
	fieldColor     = color.New()
)

func formatDecValue(value uint) string {
	return fmt.Sprintf("%d", value)
}

func formatHexValue(value uint) string {
	return fmt.Sprintf("%#08x", value)
}

// PrintHeader WDB6 header
func PrintHeader(header *Header) {
	printField("RecordCount", formatDecValue(header.RecordCount))
	printField("FieldCount", formatDecValue(header.FieldCount))
	printField("RecordSize", formatDecValue(header.RecordSize))
	printField("StringTableSize", formatDecValue(header.StringTableSize))
	printField("TableHash", formatHexValue(header.TableHash))
	printField("LayoutHash", formatHexValue(header.LayoutHash))
	printField("MinID", formatDecValue(header.MinID))
	printField("MaxID", formatDecValue(header.MaxID))
	printField("Locale", formatHexValue(header.Locale))
	printField("CopyTableSize", formatDecValue(header.CopyTableSize))
	printField("Flags", formatHexValue(header.Flags))
	printField("IDIndex", formatDecValue(header.IDIndex))
	printField("TotalFieldCount", formatDecValue(header.TotalFieldCount))
	printField("CommonDataTableSize", formatDecValue(header.CommonDataTableSize))
}

func printField(name string, value string) {
	_, _ = fieldNameColor.Printf("%s: ", name)
	_, _ = fieldColor.Printf("%s\n", value)
}
