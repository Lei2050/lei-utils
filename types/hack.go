package types

import (
	"unsafe"
)

// USE AT YOUR OWN RISK

// String force casts a []byte to a string.
func String(b []byte) (s string) {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// Bytes force casts a string to a []byte
func Bytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
