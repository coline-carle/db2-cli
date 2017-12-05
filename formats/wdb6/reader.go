package wdb6

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/wow-sweetlie/db2-cli/formats/dblayout"
)

const wdb6Magic = "WDB6"
const headerSize = 0x38

type decoder struct {
	r            io.ReadSeeker
	tmp          [256]byte
	header       *Header
	fieldsFormat []FieldFormat
	strings      map[int]string
	stringsPos   []int
	table        *dblayout.Table
	records      [][]interface{}
}

// A FormatError reports that the input is not a valid wdb6
type FormatError string

func (e FormatError) Error() string { return "WDB6: invalid format: " + string(e) }

// DecodeStrings ...
func DecodeStrings(r io.ReadSeeker) ([]int, map[int]string, error) {
	var err error

	d := &decoder{
		r: r,
	}

	if _, err = d.doDecodeHeader(); err != nil {
		return nil, nil, err
	}

	if err = d.readStrings(); err != nil {
		return nil, nil, err
	}

	return d.stringsPos, d.strings, nil
}

// DecodeFieldsFormat read db2 meta informations
func DecodeFieldsFormat(r io.ReadSeeker) ([]FieldFormat, error) {
	var err error

	d := &decoder{
		r: r,
	}
	_, err = d.doDecodeHeader()

	err = d.readFieldsFormat()

	if err != nil {
		return nil, err
	}

	return d.fieldsFormat, nil
}

// DecodeHeader ..
func DecodeHeader(r io.ReadSeeker) (*Header, error) {
	d := &decoder{
		r: r,
	}
	return d.doDecodeHeader()
}

// Decode ...
func Decode(r io.ReadSeeker, table dblayout.Table) ([][]interface{}, error) {
	var err error

	d := &decoder{
		r:     r,
		table: &table,
	}

	d.doDecodeHeader()

	if err = d.readFieldsFormat(); err != nil {
		return nil, err
	}

	if err = d.readStrings(); err != nil {
		return nil, err
	}

	if err = d.readRecords(); err != nil {
		return nil, err
	}

	return d.records, nil
}

func (d *decoder) doDecodeHeader() (h *Header, err error) {
	if err = d.checkMagic(); err != nil {
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	}

	if err = d.readHeader(); err != nil {
		return nil, err
	}

	return d.header, nil
}

func (d *decoder) decodeInt32(data []byte) interface{} {
	v := binary.LittleEndian.Uint32(data[:4])
	return int(v)
}

func (d *decoder) decodeString(data []byte) interface{} {
	v := int(binary.LittleEndian.Uint32(data[:4]))
	return d.strings[v]
}

func (d *decoder) readRecord(data []byte) (record []interface{}) {
	record = make([]interface{}, d.header.FieldCount)
	for i, fieldFormat := range d.fieldsFormat {
		pos := fieldFormat.Position
		switch d.table.Fields[i].Type {
		case "int":
			v := d.decodeInt32(data[pos:])
			record[i] = v
		case "string":
			v := d.decodeString(data[pos:])
			record[i] = v
		}
	}

	return record
}

func (d *decoder) readRecords() (err error) {
	recordBlockPosition := d.header.RecordBlockPosition()

	_, err = d.r.Seek(int64(recordBlockPosition), io.SeekStart)
	if err != nil {
		return err
	}

	len := d.header.RecordSize
	r := bufio.NewReader(d.r)
	d.records = make([][]interface{}, d.header.RecordCount)
	for i := 0; i < d.header.RecordCount; i++ {
		_, err = io.ReadFull(r, d.tmp[:len])
		d.records[i] = d.readRecord(d.tmp[:len])
	}
	return nil
}

// ReadStrings ...
func (d *decoder) readStrings() error {
	tablePos, err := d.header.StringTablePosition()
	if err != nil {
		return err
	}
	_, err = d.r.Seek(int64(tablePos), io.SeekStart)
	if err != nil {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}
		return err
	}

	b := bufio.NewReader(d.r)
	pos := 0
	str := new(bytes.Buffer)
	d.strings = make(map[int]string)
	var stringStartPos int
	stringStarted := false

	for pos < d.header.StringTableSize {
		r, runeSize, err := b.ReadRune()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}

		if r != 0x00 {
			str.WriteRune(r)
			if !stringStarted {
				stringStarted = true
				stringStartPos = pos
			}
		} else {
			if stringStarted {
				d.strings[stringStartPos] = str.String()
				d.stringsPos = append(d.stringsPos, stringStartPos)
				str = new(bytes.Buffer)
				stringStarted = false
			}
		}
		pos += runeSize
	}

	return nil
}

func (d *decoder) readFieldsFormat() (err error) {
	dataLen := d.header.FieldCount * fieldFormatSize
	_, err = io.ReadFull(d.r, d.tmp[:dataLen])
	if err != nil {
		return err
	}

	b := readBuf(d.tmp[:dataLen])
	d.fieldsFormat = make([]FieldFormat, d.header.FieldCount)

	for i := 0; i < d.header.FieldCount; i++ {
		size := int(b.uint16())
		if size > 32 {
			return fmt.Errorf("field size > 32 bits non-emplemented yet")
		}
		size = (32 - size) / 8
		position := int(b.uint16())

		fieldFormat := FieldFormat{
			Size:     size,
			Position: position,
		}
		d.fieldsFormat[i] = fieldFormat
	}

	return nil
}

func (d *decoder) readHeader() (err error) {
	_, err = io.ReadFull(d.r, d.tmp[:headerSize-len(wdb6Magic)])
	if err != nil {
		return err
	}
	b := readBuf(d.tmp[:headerSize-len(wdb6Magic)])
	d.header = &Header{
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
	return nil
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
