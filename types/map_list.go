package types

import (
	"container/list"
)

const (
	MaxNumDefault = 1000
)

type MapList[T comparable] struct {
	dataMap  map[T]*list.Element
	dataList *list.List
	maxNum   int
}

func NewMapStringList[T comparable]() *MapList[T] {
	return &MapList[T]{
		dataMap:  make(map[T]*list.Element),
		dataList: list.New(),
		maxNum:   MaxNumDefault,
	}
}

func (mapList *MapList[T]) Exists(id T) (*list.Element, bool) {
	e, exists := mapList.dataMap[id]
	return e, exists
}

func (mapList *MapList[T]) Push(key T, data interface{}) bool {
	if _, res := mapList.Exists(key); res {
		return false
	}
	elem := mapList.dataList.PushBack(data)
	mapList.dataMap[key] = elem
	return true
}

func (mapList *MapList[T]) Remove(key T) bool {
	if _, res := mapList.Exists(key); !res {
		return false
	}
	mapList.dataList.Remove(mapList.dataMap[key])
	delete(mapList.dataMap, key)
	return true
}

func (mapList *MapList[T]) Size() int {
	return mapList.dataList.Len()
}

func (mapList *MapList[T]) Pop() *list.Element {
	return mapList.dataList.Front()
}

// 遍历
func (mapList *MapList[T]) EachItem(dealFun func(interface{})) {
	for e := mapList.dataList.Front(); e != nil; e = e.Next() {
		dealFun(e.Value)
	}
}

func (mapList *MapList[T]) EachItemValue(d interface{}, dealFun func(interface{}, interface{})) {
	for e := mapList.dataList.Front(); e != nil; e = e.Next() {
		dealFun(e.Value, d)
	}
}

func (mapList *MapList[T]) EachValue() map[T]*list.Element {
	return mapList.dataMap
}

// 遍历 可根据情况中断
func (mapList *MapList[T]) EachItemBreak(dealFun func(interface{}) bool) {
	for e := mapList.dataList.Front(); e != nil; e = e.Next() {
		r := dealFun(e.Value)
		if r {
			break
		}
	}
}

func (mapList *MapList[T]) GetIndex(id uint64, dealFun func(interface{}, uint64) bool) int {
	var index int
	for e := mapList.dataList.Front(); e != nil; e = e.Next() {
		r := dealFun(e.Value, id)
		if r {
			break
		}
		index++
	}
	return index
}

func (mapList *MapList[T]) GetElement(uid []uint64, dealFun func(interface{}, []uint64) bool) *list.Element {
	for e := mapList.dataList.Front(); e != nil; e = e.Next() {
		r := dealFun(e.Value, uid)
		if r {
			return e
		}
	}

	return nil
}
