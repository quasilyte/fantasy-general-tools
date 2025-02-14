package dat

import (
	"testing"
	"unsafe"
)

func TestUnitSize(t *testing.T) {
	if s := unsafe.Sizeof(RawMagequipUnit{}); s != unitSize {
		t.Fatalf("RawMagequipUnit size is %d", s)
	}
}
