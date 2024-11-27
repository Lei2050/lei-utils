package bitmap

import (
	"bytes"
	"fmt"
	"math/bits"
)

type Bitmap []uint64

func NewBitMapBySetTo(x uint32) *Bitmap {
	return new(Bitmap).SetTo(x)
}

func (dst *Bitmap) Set(x uint32) *Bitmap {
	blkAt := int(x >> 6)
	bitAt := int(x % 64)
	if size := len(*dst); blkAt >= size {
		dst.grow(blkAt)
	}

	(*dst)[blkAt] |= (1 << bitAt)

	return dst
}

// Remove removes the bit x from the bitmap, but does not shrink it.
func (dst *Bitmap) Remove(x uint32) *Bitmap {
	if blkAt := int(x >> 6); blkAt < len(*dst) {
		bitAt := int(x % 64)
		(*dst)[blkAt] &^= (1 << bitAt)
	}
	return dst
}

func (dst *Bitmap) Clear() *Bitmap {
	//*dst = make([]uint64, 0)
	for i := 0; i < len(*dst); i++ {
		(*dst)[i] = 0
	}
	return dst
}

func (dst *Bitmap) SetTo(x uint32) *Bitmap {
	blkAt := int(x >> 6)
	if size := len(*dst); blkAt >= size {
		dst.grow(blkAt)
	}
	for i := 0; i < blkAt; i++ {
		(*dst)[i] = ^uint64(0)
	}
	bitAt := int(x % 64)
	//(*dst)[blkAt] &= (1 << bitAt)
	(*dst)[blkAt] = (1 << (bitAt + 1)) - 1
	return dst
}

// Contains checks whether a value is contained in the bitmap or not.
func (dst Bitmap) Contains(x uint32) bool {
	blkAt := int(x >> 6)
	if size := len(dst); blkAt >= size {
		return false
	}

	bitAt := int(x % 64)
	return (dst[blkAt] & (1 << bitAt)) > 0
}

// Ones sets the entire bitmap to one
func (dst Bitmap) Ones() {
	for i := 0; i < len(dst); i++ {
		dst[i] = 0xffffffffffffffff
	}
}

// Min get the smallest value stored in this bitmap, assuming the bitmap is not empty.
func (dst Bitmap) Min() (uint32, bool) {
	for blkAt, blk := range dst {
		if blk != 0x0 {
			return uint32(blkAt<<6 + bits.TrailingZeros64(blk)), true
		}
	}

	return 0, false
}

// Max get the largest value stored in this bitmap, assuming the bitmap is not empty.
func (dst Bitmap) Max() (uint32, bool) {
	var blk uint64
	for blkAt := len(dst) - 1; blkAt >= 0; blkAt-- {
		if blk = dst[blkAt]; blk != 0x0 {
			return uint32(blkAt<<6 + (63 - bits.LeadingZeros64(blk))), true
		}
	}
	return 0, false
}

// MinZero finds the first zero bit and returns its index, assuming the bitmap is not empty.
func (dst Bitmap) MinZero() (uint32, bool) {
	for blkAt, blk := range dst {
		if blk != 0xffffffffffffffff {
			return uint32(blkAt<<6 + bits.TrailingZeros64(^blk)), true
		}
	}
	return 0, false
}

// MaxZero get the last zero bit and return its index, assuming bitmap is not empty
func (dst Bitmap) MaxZero() (uint32, bool) {
	var blk uint64
	for blkAt := len(dst) - 1; blkAt >= 0; blkAt-- {
		if blk = dst[blkAt]; blk != 0xffffffffffffffff {
			return uint32(blkAt<<6 + (63 - bits.LeadingZeros64(^blk))), true
		}
	}
	return 0, false
}

// CountTo counts the number of elements in the bitmap up until the specified index. If until
// is math.MaxUint32, it will return the count. The count is non-inclusive of the index.
func (dst Bitmap) CountTo(until uint32) int {
	if len(dst) == 0 {
		return 0
	}

	// Figure out the index of the last block
	blkUntil := int(until >> 6)
	bitUntil := int(until % 64)
	if blkUntil >= len(dst) {
		blkUntil = len(dst) - 1
	}

	// Count the bits right before the last block
	sum := dst[:blkUntil].Count()

	// Count the bits at the end
	sum += bits.OnesCount64(dst[blkUntil] << (64 - uint64(bitUntil)))
	return sum
}

// Grow grows the bitmap size until we reach the desired bit.
func (dst *Bitmap) Grow(desiredBit uint32) {
	dst.grow(int(desiredBit >> 6))
}

// grow grows the size of the bitmap until we reach the desired block offset
func (dst *Bitmap) grow(blkAt int) {
	// Note that a bitmap is automatically initialized with zeros.

	// If blkAt is no greater that the current length, do nothing.
	if len(*dst) > blkAt {
		return
	}

	// If blkAt is no greater than the current capacity, resize the slice without copying.
	if cap(*dst) > blkAt {
		*dst = (*dst)[:blkAt+1]
		return
	}

	old := *dst
	*dst = make(Bitmap, blkAt+1, capacityFor(blkAt+1))
	copy(*dst, old)
}

// balance grows the destination bitmap to match the size of the source bitmap.
func (dst *Bitmap) balance(src Bitmap) {
	if len(*dst) < len(src) {
		dst.grow(len(src) - 1)
	}
}

// capacityFor computes the next power of 2 for a given index
func capacityFor(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return int(v)
}

// And computes the intersection between two bitmaps and stores the result in the current bitmap
func (dst *Bitmap) And(b Bitmap) {
	if dst.balance(b); len(*dst) >= len(b) {
		and(dst, b)
	}
}

// AndNot computes the difference between two bitmaps and stores the result in the current bitmap
func (dst *Bitmap) AndNot(b Bitmap) {
	if dst.balance(b); len(*dst) >= len(b) {
		andn(dst, b)
	}
}

// Or computes the union between two bitmaps and stores the result in the current bitmap
func (dst *Bitmap) Or(b Bitmap) {
	if dst.balance(b); len(*dst) >= len(b) {
		or(dst, b)
	}
}

// Xor computes the symmetric difference between two bitmaps and stores the result in the current bitmap
func (dst *Bitmap) Xor(b Bitmap) {
	if dst.balance(b); len(*dst) >= len(b) {
		xor(dst, b)
	}
}

// Count returns the number of elements in this bitmap
func (dst Bitmap) Count() int {
	return count(dst)
}

func (dst Bitmap) Dump() {
	dump(dst)
}

// Count counts the number of bits set to one
func count(arr []uint64) int {
	sum := 0
	for i := 0; i < len(arr); i++ {
		sum += bits.OnesCount64(arr[i])
	}
	return sum
}

func (dst Bitmap) IsEqual(b Bitmap) bool {
	ld, lb := len(dst), len(b)
	for i := 0; i < ld && i < lb; i++ {
		if dst[i] != b[i] {
			return false
		}
	}

	var (
		tmp *Bitmap = &dst
		s   int     = lb
	)
	if ld < lb {
		tmp, s = &b, ld
	}
	for i := s; i < len(*tmp); i++ {
		if (*tmp)[i] != 0 {
			return false
		}
	}
	return true
}

// and computes the intersection between two bitmaps and stores the result in the current bitmap
func and(dst *Bitmap, b Bitmap) {
	a := *dst
	for i := 0; i < len(b); i++ {
		a[i] = a[i] & b[i]
	}
}

// AndNot computes the difference between two bitmaps and stores the result in the current bitmap
func andn(dst *Bitmap, b Bitmap) {
	a := *dst
	for i := 0; i < len(b); i++ {
		a[i] = a[i] &^ b[i]
	}
}

// or computes the union between two bitmaps and stores the result in the current bitmap
func or(dst *Bitmap, b Bitmap) {
	a := *dst
	for i := 0; i < len(b); i++ {
		a[i] = a[i] | b[i]
	}
}

// Xor computes the symmetric difference between two bitmaps and stores the result in the current bitmap
func xor(dst *Bitmap, b Bitmap) {
	a := *dst
	for i := 0; i < len(b); i++ {
		a[i] = a[i] ^ b[i]
	}
}

func dump(b Bitmap) {
	var buffer bytes.Buffer
	var scale bytes.Buffer
	var m int
	for i := 0; i < len(b); i++ {
		for j := 0; j < 64; j++ {
			if (1<<j)&b[i] > 0 {
				buffer.WriteString("1")
			} else {
				buffer.WriteString("0")
			}

			if m%64 == 0 {
				scale.WriteString("|")
			} else {
				scale.WriteString("-")
			}
			m++
		}
		//buffer.WriteString(fmt.Sprintf("%064b", b[i]))
	}
	fmt.Println(scale.String())
	fmt.Println(buffer.String())
}
