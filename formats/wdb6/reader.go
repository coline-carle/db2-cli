package wdb6

import (
	"encoding/binary"
	"io"
)

const wdb6Magic = "WDB6"
const headerSize = 0x34

type decoder struct {
	r   io.Reader
	tmp [256]byte
}

// A FormatError reports that the input is not a valid wdb6
type FormatError string

func (e FormatError) Error() string { return "WDB6: invalid format: " + string(e) }

// Decode a db2 file with WDB6 format
func Decode(r io.Reader) (*Header, error) {
	d := &decoder{
		r: r,
	}

	if err := d.checkMagic(); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	header, err := d.readHeader()
	return header, err
}

func (d *decoder) readHeader() (header *Header, err error) {
	_, err = io.ReadFull(d.r, d.tmp[:headerSize])
	if err != nil {
		return nil, err
	}
	b := readBuf(d.tmp[:headerSize])
	h := &Header{
		RecordCount:         uint(b.uint32()),
		FieldCount:          uint(b.uint32()),
		RecordSize:          uint(b.uint32()),
		StringTableSize:     uint(b.uint32()),
		TableHash:           uint(b.uint32()),
		LayoutHash:          uint(b.uint32()),
		MinID:               uint(b.uint32()),
		MaxID:               uint(b.uint32()),
		Locale:              uint(b.uint32()),
		CopyTableSize:       uint(b.uint32()),
		Flags:               uint(b.uint16()),
		IDIndex:             uint(b.uint16()),
		TotalFieldCount:     uint(b.uint32()),
		CommonDataTableSize: uint(b.uint32()),
	}
	return h, nil
}

func (d *decoder) checkMagic() error {
	_, err := io.ReadFull(d.r, d.tmp[:len(wdb6Magic)])
	if err != nil {
		return err
	}
	if string(d.tmp[:len(wdb6Magic)]) != wdb6Magic {
		return FormatError("not a DB6 file")
	}
	return nil
}

type readBuf []byte

func (b *readBuf) uint8() uint8 {
	v := (*b)[0]
	*b = (*b)[1:]
	return v
}

func (b *readBuf) uint16() uint16 {
	v := binary.LittleEndian.Uint16(*b)
	*b = (*b)[2:]
	return v
}

func (b *readBuf) uint32() uint32 {
	v := binary.LittleEndian.Uint32(*b)
	*b = (*b)[4:]
	return v
}
