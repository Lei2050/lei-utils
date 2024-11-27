package types

import (
	"fmt"
	"testing"
)

func TestHack1(t *testing.T) {
	str := "hello world"
	slc := Bytes(str)
	fmt.Println(str, slc)

	//slc[1] = 'A' //panic: unexpected fault address
	//fmt.Println(str, slc)
}

func TestHack2(t *testing.T) {
	slc := []byte("hello world")
	str := String(slc)
	t.Log(str, slc)

	slc[1] = 'A' //panic: unexpected fault address
	t.Log(str, slc)
}
