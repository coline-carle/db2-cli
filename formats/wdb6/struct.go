package wdb6

// Wdb6 file format
type Wdb6 struct {
	Header       *Header
	FieldsFormat []FieldFormat
}

const (
	flagHasOffsetMap    = 0x01
	flagHasSecondaryKey = 0x02
	flagHasNonInlineID  = 0x04
)

// Header of Wdb6 file
type Header struct {
	RecordCount         uint
	FieldCount          uint
	RecordSize          uint
	StringTableSize     uint
	TableHash           uint
	LayoutHash          uint
	MinID               uint
	MaxID               uint
	Locale              uint
	CopyTableSize       uint
	Flags               uint
	IDIndex             uint
	TotalFieldCount     uint
	CommonDataTableSize uint
}

// FieldFormat description
type FieldFormat struct {
	Size     uint
	Position uint
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
