package dat

import (
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

func scanUint8(data *[]byte) (v uint8) {
	v = (*data)[0]
	*data = (*data)[1:]
	return v
}

func scanUint16(data *[]byte) (v uint16) {
	v = binary.LittleEndian.Uint16(*data)
	*data = (*data)[2:]
	return v
}

func scanUint32(data *[]byte) (v uint32) {
	v = binary.LittleEndian.Uint32(*data)
	*data = (*data)[4:]
	return v
}

func scanUint32BE(data *[]byte) (v uint32) {
	v = binary.BigEndian.Uint32(*data)
	*data = (*data)[4:]
	return v
}
