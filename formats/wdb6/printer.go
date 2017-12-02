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
	titleColor     = color.New(color.FgRed).Add(color.Bold)
	chapterColor   = color.New(color.FgGreen).Add(color.Underline)
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
	if len(flags) > 0 {
		_, _ = noteColor.Printf("(%s)", strings.Join(flags, ", "))
	}

	fmt.Printf("\n")
}

func printFieldSize(size uint) {
	fieldNameColor.Print("Size: ")
	fieldColor.Print(fmt.Sprintf("%d B\n", size))
}

// PrintFieldsFormat print field formats
func PrintFieldsFormat(fields []FieldFormat) {
	chapterColor.Println("WDB6 fields format")
	for index, field := range fields {
		title := fmt.Sprintf("Field %d", index)
		titleColor.Println(title)

		printFieldSize(field.Size)
		hexPos := fmt.Sprintf("%#04x", field.Position)
		printField("Position", hexPos)
		fmt.Println()
	}
}

// PrintHeader WDB6 header
func PrintHeader(header *Header) {
	chapterColor.Println("WDB6 Header")
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
	printFlags(header)
	printField("IDIndex", formatDecValue(header.IDIndex))
	printField("TotalFieldCount", formatDecValue(header.TotalFieldCount))
	printField("CommonDataTableSize", formatDecValue(header.CommonDataTableSize))
}

func printField(name string, value string) {
	_, _ = fieldNameColor.Printf("%s: ", name)
	_, _ = fieldColor.Printf("%s\n", value)
}
