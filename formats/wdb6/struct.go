package wdb6

import "fmt"

const (
	flagHasOffsetMap    = 0x01
	flagHasSecondaryKey = 0x02
	flagHasNonInlineID  = 0x04
)

var (
	errNoStringTable     = fmt.Errorf("no string table present")
	errNoOffsetMap       = fmt.Errorf("no offset map present")
	errNoCopyTable       = fmt.Errorf("no copy table present")
	errNoCommonDataTable = fmt.Errorf("no common data table present")
)

const fieldFormatSize = 0x4

// Header of Wdb6 file
type Header struct {
	RecordCount         int
	FieldCount          int
	RecordSize          int
	StringTableSize     int
	TableHash           int
	LayoutHash          int
	MinID               int
	MaxID               int
	Locale              int
	CopyTableSize       int
	Flags               int
	IDIndex             int
	TotalFieldCount     int
	CommonDataTableSize int
}

// FieldFormat description
type FieldFormat struct {
	Size     int
	Position int
}

// HasStringTable return true if the db2 file contain strings
func (h *Header) HasStringTable() bool {
	return !h.HasOffsetMap()
}

// HasOffsetMap ...
func (h *Header) HasOffsetMap() bool {
	return (h.Flags & flagHasOffsetMap) == flagHasOffsetMap
}

// HasSecondaryKey ...
func (h *Header) HasSecondaryKey() bool {
	return (h.Flags & flagHasSecondaryKey) == flagHasSecondaryKey
}

// HasNonInlineID ...
func (h *Header) HasNonInlineID() bool {
	return (h.Flags & flagHasNonInlineID) == flagHasNonInlineID
}

// RecordBlockPosition ...
func (h *Header) RecordBlockPosition() int {
	return headerSize + h.FieldFormatBlockSize()
}

// RecordBlockSize ....
func (h *Header) RecordBlockSize() int {
	return h.RecordSize * h.RecordCount
}

// StringTablePosition ...
func (h *Header) StringTablePosition() (int, error) {
	if !h.HasStringTable() {
		return 0, errNoStringTable
	}

	return h.RecordBlockPosition() + h.RecordBlockSize(), nil
}

// OffsetMapPosition ...
func (h *Header) OffsetMapPosition() (int, error) {
	if !h.HasOffsetMap() {
		return 0, errNoOffsetMap
	}
	position := h.RecordBlockPosition()
	position += h.RecordBlockSize()
	position += h.StringTableSize
	return position, nil
}

// FieldFormatBlockSize ...
func (h *Header) FieldFormatBlockSize() int {
	return h.FieldCount * fieldFormatSize
}
