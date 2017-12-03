package wdb6

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const wdb6Magic = "WDB6"
const headerSize = 0x38

type decoder struct {
	r   io.Reader
	tmp [256]byte
}

// A FormatError reports that the input is not a valid wdb6
type FormatError string

func (e FormatError) Error() string { return "WDB6: invalid format: " + string(e) }

// Decode a db2 file with WDB6 format
func Decode(r io.Reader) (*Wdb6, error) {
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
	if err != nil {
		return nil, err
	}
	fieldsFormat, err := d.readFieldsFormat(header)
	if err != nil {
		return nil, err
	}
	wdb6 := &Wdb6{
		Header:       header,
		FieldsFormat: fieldsFormat,
	}

	return wdb6, nil
}

// ReadStrings ...
func ReadStrings(db2 *Wdb6, reader io.ReadSeeker) ([]int, map[int]string, error) {
	tablePos, err := db2.Header.StringTablePosition()
	if err != nil {
		return nil, nil, err
	}
	_, err = reader.Seek(int64(tablePos), io.SeekStart)
	if err != nil {
		if err == io.EOF {
			return nil, nil, fmt.Errorf("unexpected EOF")
		}
		return nil, nil, err
	}

	b := bufio.NewReader(reader)
	pos := 0
	str := new(bytes.Buffer)
	strings := make(map[int]string)
	var positions []int
	var stringStartPos int
	stringStarted := false

	for pos < db2.Header.StringTableSize {
		r, runeSize, err := b.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, nil, fmt.Errorf("unexpected EOF")
			}
			return nil, nil, err
		}

		if r != 0x00 {
			str.WriteRune(r)
			if !stringStarted {
				stringStarted = true
				stringStartPos = pos
			}
		} else {
			if stringStarted {
				strings[stringStartPos] = str.String()
				positions = append(positions, stringStartPos)
				str = new(bytes.Buffer)
				stringStarted = false
			}
		}
		pos += runeSize
	}

	return positions, strings, nil
}

func (d *decoder) readFieldsFormat(header *Header) (fieldFormat []FieldFormat, err error) {
	dataLen := header.FieldCount * fieldFormatSize
	_, err = io.ReadFull(d.r, d.tmp[:dataLen])
	if err != nil {
		return []FieldFormat{}, err
	}

	b := readBuf(d.tmp[:dataLen])
	fieldsFormat := make([]FieldFormat, header.FieldCount)

	for i := 0; i < header.FieldCount; i++ {
		size := int(b.uint16())
		if size > 32 {
			return []FieldFormat{}, fmt.Errorf("field size > 32 bits non-emplemented yet")
		}
		size = (32 - size) / 8
		position := int(b.uint16())

		fieldFormat := FieldFormat{
			Size:     size,
			Position: position,
		}
		fieldsFormat[i] = fieldFormat
	}

	return fieldsFormat, nil
}

func (d *decoder) readHeader() (header *Header, err error) {
	_, err = io.ReadFull(d.r, d.tmp[:headerSize-len(wdb6Magic)])
	if err != nil {
		return nil, err
	}
	b := readBuf(d.tmp[:headerSize-len(wdb6Magic)])
	h := &Header{
		RecordCount:         int(b.uint32()),
		FieldCount:          int(b.uint32()),
		RecordSize:          int(b.uint32()),
		StringTableSize:     int(b.uint32()),
		TableHash:           int(b.uint32()),
		LayoutHash:          int(b.uint32()),
		MinID:               int(b.uint32()),
		MaxID:               int(b.uint32()),
		Locale:              int(b.uint32()),
		CopyTableSize:       int(b.uint32()),
		Flags:               int(b.uint16()),
		IDIndex:             int(b.uint16()),
		TotalFieldCount:     int(b.uint32()),
		CommonDataTableSize: int(b.uint32()),
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
