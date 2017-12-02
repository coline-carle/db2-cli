package wdb6

// Wdb6 file format
type Wdb6 struct {
	Header Header
}

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
