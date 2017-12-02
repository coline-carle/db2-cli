package wdb6

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

var (
	fieldNameColor = color.New(color.FgCyan).Add(color.Bold)
	fieldColor     = color.New()
	noteColor      = color.New(color.FgYellow)
)

func formatDecValue(value uint) string {
	return fmt.Sprintf("%d", value)
}

func formatHexValue(value uint) string {
	return fmt.Sprintf("%#08x", value)
}

func printFlags(h *Header) {
	var flags []string
	hexValue := formatHexValue(h.Flags)

	if h.HasOffsetMap() {
		flags = append(flags, "offset map")
	}

	if h.HasSecondaryKey() {
		flags = append(flags, "secondary key")
	}

	if h.HasNonInlineID() {
		flags = append(flags, "non inline ID")
	}

	_, _ = fieldNameColor.Print("Flags: ")
	_, _ = fieldColor.Printf("%s  ", hexValue)
	_, _ = noteColor.Printf("(%s)\n", strings.Join(flags, ", "))
}

// Print WDB6 header
func (header *Header) Print() {
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
	// printField("Flags", formatHexValue(header.Flags))
	printFlags(header)
	printField("IDIndex", formatDecValue(header.IDIndex))
	printField("TotalFieldCount", formatDecValue(header.TotalFieldCount))
	printField("CommonDataTableSize", formatDecValue(header.CommonDataTableSize))
}

func printField(name string, value string) {
	_, _ = fieldNameColor.Printf("%s: ", name)
	_, _ = fieldColor.Printf("%s\n", value)
}
