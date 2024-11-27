package bitmap

import (
	"fmt"
	"testing"
)

func TestEventMgr(t *testing.T) {
	var b Bitmap
	b.Set(1).Set(2).Set(3).Set(64)
	b.Dump()

	fmt.Printf("1=%v, 2=%v, 4=%v, 63=%v, 64=%v, 65=%v\n", b.Contains(1), b.Contains(2), b.Contains(4), b.Contains(63), b.Contains(64), b.Contains(65))

	b.Remove(1).Remove(3)
	b.Dump()

	fmt.Printf("1=%v, 2=%v, 4=%v, 63=%v, 64=%v, 65=%v\n", b.Contains(1), b.Contains(2), b.Contains(4), b.Contains(63), b.Contains(64), b.Contains(65))

	b.Clear()
	b.Set(64).Set(128)
	b.Dump()

	b.Clear().SetTo(5)
	b.Dump()

	b.Clear().SetTo(64)
	b.Dump()

	b.Clear().SetTo(256 + 1)
	b.Dump()
	b.Clear()

	var a Bitmap
	a.Set(1)
	b.Set(2)
	fmt.Printf("%+v\n", a.IsEqual(b))

	a.Clear().Set(66)
	b.Clear().Set(66)
	fmt.Printf("%+v\n", a.IsEqual(b))

	a.Clear().SetTo(66)
	b.Clear().Set(1024)
	fmt.Printf("%+v\n", a.IsEqual(b))
	fmt.Printf("%+v\n", b.IsEqual(a))

	a.Clear().Set(1024)
	b.Clear().Set(1024)
	fmt.Printf("%+v\n", a.IsEqual(b))
	fmt.Printf("%+v\n", b.IsEqual(a))
}
