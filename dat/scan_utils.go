package dat

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

// uint16le is a uint16-like value from the serialized format.
// Encoded as [low][high], 19 is [19][0] (Little Endian).
type uint16le struct {
	Low  byte
	High byte
}

func (d uint16le) ToInt() int {
	return int(d.Low) | (int(d.High) << 8)
}

func (d *uint16le) SetInt(v int) {
	d.Low = byte(v)
	d.High = byte(v >> 8)
}

func asBytes[T any](obj *T) []byte {
	length := unsafe.Sizeof(*obj)
	return unsafe.Slice((*byte)(unsafe.Pointer(obj)), length)
}

func trimCstring(raw []byte) string {
	nullByte := bytes.IndexByte(raw, 0)
	if nullByte != -1 {
		raw = raw[:nullByte]
	}
	return string(raw)
}

func scanString(data *[]byte, length int) string {
	realLength := length
	nullByte := bytes.IndexByte(*data, 0)
	if nullByte != -1 && nullByte < length {
		realLength = nullByte
	}

	s := string((*data)[:realLength])
	*data = (*data)[length:]
	return s
}

func scanBool(data *[]byte) bool {
	v := (*data)[0]
	*data = (*data)[1:]
	return v != 0
}

func scanUint8(data *[]byte) uint8 {
	v := (*data)[0]
	*data = (*data)[1:]
	return v
}

func scanUint16(data *[]byte) uint16 {
	v := binary.LittleEndian.Uint16(*data)
	*data = (*data)[2:]
	return v
}

func scanUint32(data *[]byte) uint32 {
	v := binary.LittleEndian.Uint32(*data)
	*data = (*data)[4:]
	return v
}

func scanUint32BE(data *[]byte) uint32 {
	v := binary.BigEndian.Uint32(*data)
	*data = (*data)[4:]
	return v
}
