package types

func BoolToNumber[T Number](b bool) T {
	if b {
		return 1
	}
	return 0
}

func TernaryOperator[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}
