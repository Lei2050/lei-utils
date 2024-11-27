package crontab

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Lei2050/lei-utils/bitmap"
	"github.com/Lei2050/lei-utils/types"
)

type Entry struct {
	minute bitmap.Bitmap //分
	hour   bitmap.Bitmap //时
	dom    bitmap.Bitmap //月天
	month  bitmap.Bitmap //月
	dow    bitmap.Bitmap //周天

	id   uint64
	cmd  func(args ...interface{})
	args []interface{}
}

func newEntry(id uint64,
	minute, hour, dom, month, dow string,
	cmd func(...interface{}), args ...interface{}) (*Entry, error) {
	e := &Entry{
		id:   id,
		cmd:  cmd,
		args: args,
	}

	var err error

	err = getRange(&e.minute, FirstMinute, LastMinute, minute)
	if err != nil {
		return nil, err
	}

	err = getRange(&e.hour, FirstHour, LastHour, hour)
	if err != nil {
		return nil, err
	}

	err = getRange(&e.dom, FirstDOM, LastDOM, dom)
	if err != nil {
		return nil, err
	}

	err = getRange(&e.month, FirstMonth, LastMonth, month)
	if err != nil {
		return nil, err
	}

	err = getRange(&e.dow, FirstDOW, LastDOW, dow)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func getRange(bitset *bitmap.Bitmap, low, high int, input string) error {
	var num1, num2, num3 int
	if input == "" || input == "*" {
		num1, num2, num3 = low, high, 1
	} else if strings.Contains(input, ",") {
		nums := types.StringToIntegerSlice[uint32](input, ",")
		if nums == nil {
			return fmt.Errorf("parse numbers failed:%s", input)
		}
		for _, v := range nums {
			bitset.Set(v)
		}
		return nil
	} else if strings.Contains(input, "-") {
		nums := types.StringToIntegerSlice[int](input, "-")
		if nums == nil {
			return fmt.Errorf("parse range numbers failed:%s", input)
		}
		if len(nums) != 2 {
			return fmt.Errorf("parse range numbers failed:%s", input)
		}
		num1, num2, num3 = nums[0], nums[1], 1
	} else {
		num, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return err
		}
		num1, num2, num3 = int(num), int(num), 1
	}

	for i := num1; i <= num2; i += num3 {
		bitset.Set(uint32(i))
	}

	return nil
}

func (e *Entry) Test(t time.Time) bool {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return e.minute.Contains(uint32(t.Minute())) &&
		e.hour.Contains(uint32(t.Hour())) &&
		e.dom.Contains(uint32(t.Day())) &&
		e.month.Contains(uint32(t.Month())) &&
		e.dow.Contains(uint32(weekday))
}
