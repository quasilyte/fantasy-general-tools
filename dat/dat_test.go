package dat

import (
	"testing"
	"unsafe"
)

func TestUnitSize(t *testing.T) {
	if s := unsafe.Sizeof(MagequipUnit{}); s != unitSize {
		t.Fatalf("MagequipUnit size is %d", s)
	}
}
