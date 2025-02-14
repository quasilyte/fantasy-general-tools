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

func scanUint8(data []byte) (v uint8, rest []byte) {
	v = data[0]
	return v, data[1:]
}

func scanUint16(data []byte) (v uint16, rest []byte) {
	v = binary.LittleEndian.Uint16(data)
	return v, data[2:]
}

func scanUint32BE(data []byte) (v uint32, rest []byte) {
	v = binary.BigEndian.Uint32(data)
	return v, data[4:]
}
