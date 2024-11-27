package types

func SliceToAnys[T any](arr []T) []any {
	ret := make([]any, len(arr))
	if len(arr) <= 0 {
		return ret
	}
	for k, v := range arr {
		ret[k] = v
	}
	return ret
}
