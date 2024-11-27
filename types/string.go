package types

import (
	"fmt"
	"strconv"
	"strings"
)

func NumberToString[T Number](v T) string {
	return fmt.Sprintf("%v", v)
}

func IntegerToString[T Integer](v T) string {
	return fmt.Sprintf("%d", v)
}

func FloatToString[T Float](v T) string {
	return fmt.Sprintf("%f", v)
}

func StringToNumber[T Number](s string) T {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return T(v)
}

func StringToInteger[T Integer](s string) T {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return T(v)
}

func StringToFloat[T Number](s string) T {
	return StringToNumber[T](s)
}

func StringToNumberSlice[T Number](s, sep string) []T {
	str1 := strings.Split(s, sep)
	if len(str1) == 0 {
		return nil
	}
	ret := make([]T, len(str1))
	for i, v := range str1 {
		ret[i] = StringToNumber[T](v)
	}
	return ret
}

func StringToIntegerSlice[T Integer](s, sep string) []T {
	str1 := strings.Split(s, sep)
	if len(str1) == 0 {
		return nil
	}
	ret := make([]T, len(str1))
	for i, v := range str1 {
		ret[i] = StringToInteger[T](v)
	}
	return ret
}

func StringToFloatSlice[T Float](s, sep string) []T {
	return StringToNumberSlice[T](s, sep)
}

// ----- slice to str
func NumberSliceToString[T Number](slc []T, sep string) string {
	var s string
	for k, v := range slc {
		s += NumberToString(v)
		if k != len(slc)-1 {
			s += sep
		}
	}
	return s
}

func IntegerSliceToString[T Integer](slc []T, sep string) string {
	var s string
	for k, v := range slc {
		s += IntegerToString(v)
		if k != len(slc)-1 {
			s += sep
		}
	}
	return s
}

func FloatSliceToString[T Float](slc []T, sep string) string {
	return NumberSliceToString(slc, sep)
}

// ------ map to str
func MapToString[KEY comparable, VALUE any](m map[KEY]VALUE, sep1, sep2 string) string {
	var s string
	first := true
	for k, v := range m {
		if first {
			s += fmt.Sprintf("%v", k) + sep1 + fmt.Sprintf("%v", v)
		} else {
			s += sep2 + fmt.Sprintf("%v", k) + sep1 + fmt.Sprintf("%v", v)
			first = false
		}
	}
	return s
}

func StringToNumberMap[KEY, VALUE Number](s string, sep1, sep2 string) map[KEY]VALUE {
	str1 := strings.Split(s, sep1)
	if len(str1) == 0 {
		return nil
	}

	m := make(map[KEY]VALUE)
	for _, v := range str1 {
		str2 := strings.Split(v, sep2)
		if len(str2) != 2 {
			continue
		}

		m[StringToNumber[KEY](str2[0])] = StringToNumber[VALUE](str2[1])
	}
	return m
}

func StringMapToNumberMap[KEY, VALUE Number](m map[string]string) map[KEY]VALUE {
	if len(m) <= 0 {
		return nil
	}

	ret := make(map[KEY]VALUE)
	for k, v := range m {
		ku := StringToNumber[KEY](k)
		vu := StringToNumber[VALUE](v)
		ret[ku] = vu
	}
	return ret
}

// gm 转换格式使用
func StringToAwardSlice(s string, sep1, sep2 string) (ts []uint32, vs []uint64, ss []uint32) {
	s1 := strings.Split(s, sep1)
	for _, s2 := range s1 {
		s3 := strings.Split(s2, sep2)
		if len(s3) != 3 {
			continue
		}
		ts = append(ts, StringToInteger[uint32](s3[0]))
		vs = append(vs, StringToInteger[uint64](s3[1]))
		ss = append(ss, StringToInteger[uint32](s3[2]))
	}
	return
}

func TransToCrontabSchedule(min, hour, day, month, week string) string {
	if min == "" {
		min = "*"
	}
	if hour == "" {
		hour = "*"
	}
	if day == "" {
		day = "*"
	}
	if month == "" {
		month = "*"
	}
	if week == "" {
		week = "*"
	}

	return min + " " + hour + " " + day + " " + month + " " + week
}
