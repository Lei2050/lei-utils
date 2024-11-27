package types

// comparable，官方已经定义了所有可用 != 以及 == 对比的类型

// Ordered 代表所有可比大小排序的类型
type Ordered interface {
	Integer | Float | ~string
}

// 数值类型
type Number interface {
	Integer | Float
}

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}
