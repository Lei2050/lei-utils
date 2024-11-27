package types

import (
	"sync"
)

type VectorFun func(any, any) bool

type Op func(any)

type Vector struct {
	Values []any
	len    int
	equal  VectorFun
	sync.RWMutex
}

func NewVector(vfunc VectorFun) *Vector {
	return &Vector{len: 0, equal: vfunc}
}

func (v *Vector) Len() int {
	v.RLock()
	defer v.RUnlock()
	return v.len
}

func (v *Vector) Index(i int) any {
	v.RLock()
	defer v.RUnlock()
	if i >= v.len {
		return nil
	}
	return v.Values[i]
}

func (v *Vector) expand() {
	curcap := len(v.Values)
	var newcap int
	if curcap == 0 {
		newcap = 8
	} else if curcap < 1024 {
		newcap = curcap * 2
	} else {
		newcap = curcap + (curcap / 4)
	}
	values := make([]any, newcap)
	if curcap != 0 {
		copy(values, v.Values)
	}
	v.Values = values
}

func (v *Vector) PushBack(value any) {
	v.Lock()
	defer v.Unlock()
	if len(v.Values) == v.len {
		v.expand()
	}
	v.Values[v.len] = value
	v.len++
}

func (v *Vector) PopBack() any {
	v.Lock()
	defer v.Unlock()
	if v.len == 0 {
		return nil
	}
	v.len--
	return v.Values[v.len]
}

func (v *Vector) Remove(value any) {
	v.Lock()
	defer v.Unlock()
	for i := 0; i < v.len; {
		if v.equal(v.Values[i], value) {
			if tmp := i + 1; tmp < v.len {
				copy(v.Values[i:], v.Values[tmp:v.len])
			}
			v.len--
		} else {
			i++
		}
	}
}

func (v *Vector) Traverse(op Op) {
	v.Lock()
	defer v.Unlock()
	for i := 0; i < v.len; i++ {
		op(v.Values[i])
	}
}
