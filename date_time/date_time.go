package date_time

import (
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	FirstDay      int32    = 2361222 //对应 1752-09-14 这天
	SecsPerDay    int32    = 86400   //<一天的时间>
	HalfPerDay    int32    = 43200   //半天時間
	MsecsPerDay   int32    = 86400000
	SecsPerHour   int32    = 3600
	MsecsPerHour  int32    = 3600000
	SecsPerMin    int32    = 60
	MsecsPerMin   int32    = 60000
	MsecsPerSec   int32    = 1000
	FirstYear     int32    = 1752
	anDaysInMonth []int32  = []int32{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	astrWeeks     []string = []string{"", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
)

type TDate struct {
	nDateVal int32
}

// 获取当前日期
func CurDate() TDate {
	nYear, nMonth, nDay := time.Now().Date()
	val := gregToJulian(int32(nYear), int32(nMonth), int32(nDay))
	return TDate{val}
}

func NewDefaultDate() *TDate {
	return &TDate{0}
}

// 创建指定日期
func NewSpecDate(nYear, nMonth, nDay int32) (*TDate, error) {
	tRetDate := &TDate{0}

	if !tRetDate.SetDate(nYear, nMonth, nDay) {
		tRetDate = nil
		return nil, errors.New("invalid param")
	}

	return tRetDate, nil
}

// 按照%04d-%02d-%02d格式创建字符串指定日期
func NewStringDate(strDate string) (*TDate, error) {
	var nYear, nMonth, nDay int32

	if _, err := fmt.Sscanf(strDate, "%04d-%02d-%02d", &nYear, &nMonth, &nDay); err != nil {
		return nil, err
	}

	tRetDate := &TDate{0}
	if !tRetDate.SetDate(nYear, nMonth, nDay) {
		tRetDate = nil
		return nil, errors.New("invalid date")
	}

	return tRetDate, nil
}

func (td TDate) IsNull() bool {
	return td.nDateVal == 0
}

func (td TDate) IsValid() bool {
	return (td.nDateVal >= FirstDay)
}

func gregToJulian(nYear, nMonth, nDay int32) int32 {
	var nC, nYa int32

	if nYear <= 99 {
		nYear += 1900
	}

	if nMonth > 2 {
		nMonth -= 3
	} else {
		nMonth += 9
		nYear--
	}

	nC = nYear
	nC /= 100
	nYa = nYear - 100*nC

	return 1721119 + nDay + ((146097 * nC) / 4) + ((1461 * nYa) / 4) + ((153*nMonth + 2) / 5)
}

func julianToGreg(nJulian int32, pnYear, pnMonth, pnDay *int32) {
	var nX int32
	var nJ int32 = nJulian - 1721119

	*pnYear = ((nJ * 4) - 1) / 146097
	nJ = (nJ * 4) - (146097 * *pnYear) - 1
	nX = nJ / 4
	nJ = ((nX * 4) + 3) / 1461
	*pnYear = (100 * *pnYear) + nJ
	nX = (nX * 4) + 3 - (1461 * nJ)
	nX = (nX + 4) / 4
	*pnMonth = ((5 * nX) - 3) / 153
	nX = (5 * nX) - 3 - (153 * *pnMonth)
	*pnDay = (nX + 5) / 5

	if *pnMonth < 10 {
		*pnMonth += 3
	} else {
		*pnMonth -= 9
		*pnYear++
	}
}

func (td TDate) Year() int32 {
	var nYear, nMonth, nDay int32
	julianToGreg(td.nDateVal, &nYear, &nMonth, &nDay)
	return nYear
}

func (td TDate) Month() int32 {
	var nYear, nMonth, nDay int32
	julianToGreg(td.nDateVal, &nYear, &nMonth, &nDay)
	return nMonth
}

func (td TDate) Day() int32 {
	var nYear, nMonth, nDay int32
	julianToGreg(td.nDateVal, &nYear, &nMonth, &nDay)
	return nDay
}

// 获取目前日期对应周几，返回(1-7)
func (td TDate) DayOfWeek() int32 {
	return (((td.nDateVal+1)%7)+6)%7 + 1
}

// 获取目前日期是所在年份中的第多少天(比如元旦表示第一天)
func (td TDate) DayOfYear() int32 {
	return td.nDateVal - gregToJulian(td.Year(), 1, 1) + 1
}

// 判断指定年份是否是闰年
func (td TDate) IsLeapYear(nYear int32) bool {
	return (nYear%4 == 0 && nYear%100 != 0 || nYear%400 == 0)
}

// 获取目前日期所在月份的总共天数
func (td TDate) DaysInMonth() int32 {
	var nYear, nMonth, nDay int32
	julianToGreg(td.nDateVal, &nYear, &nMonth, &nDay)
	if nMonth == 2 && td.IsLeapYear(nYear) {
		return 29
	} else {
		return anDaysInMonth[nMonth]
	}
}

// 获取目前日期所在年份的总共天数
func (td TDate) DaysInYear() int32 {
	var nYear, nMonth, nDay int32
	julianToGreg(td.nDateVal, &nYear, &nMonth, &nDay)
	if td.IsLeapYear(nYear) {
		return 366
	} else {
		return 365
	}
}

// 获取目前日期是周几的字符串
func (td TDate) DayOfWeekToString() string {
	return astrWeeks[td.DayOfWeek()]
}

func (td TDate) isValid(nYear, nMonth, nDay int32) bool {
	if nYear >= 0 && nYear <= 99 {
		nYear += 1900
	}

	if nYear < FirstYear || (nYear == FirstYear && (nMonth < 9 || (nMonth == 9 && nDay < 14))) {
		return false
	}

	return (nDay > 0 && nMonth > 0 && nMonth <= 12) && (nDay <= anDaysInMonth[nMonth] || (nDay == 29 && nMonth == 2 && td.IsLeapYear(nYear)))
}

// 将目前日志设置为指定日期
func (td *TDate) SetDate(nYear, nMonth, nDay int32) bool {
	if !td.isValid(nYear, nMonth, nDay) {
		return false
	}

	td.nDateVal = gregToJulian(nYear, nMonth, nDay)
	return true
}

// 返回目前日期加上指定天数对应的日期，参数可正可负
func (td TDate) AddDays(nDays int32) TDate {
	return TDate{td.nDateVal + nDays}
}

// 返回目前日期距离指定日期的天数差，为负值时表示目前日期相比指定日期是未来时间
func (td TDate) DaysTo(ptDstDate *TDate) int32 {
	return ptDstDate.nDateVal - td.nDateVal
}

// 将目前日期按照%04d-%02d-%02d格式转换为字符串
func (td TDate) ToString() string {
	var nYear, nMonth, nDay int32
	julianToGreg(td.nDateVal, &nYear, &nMonth, &nDay)

	return fmt.Sprintf("%04d-%02d-%02d", nYear, nMonth, nDay)
}

// 以下就是一些日期比较函数
func (td TDate) IsEqual(ptDate *TDate) bool {
	return td.nDateVal == ptDate.nDateVal
}

func (td TDate) IsNotEqual(ptDate *TDate) bool {
	return td.nDateVal != ptDate.nDateVal
}

func (td TDate) IsLess(ptDate *TDate) bool {
	return td.nDateVal < ptDate.nDateVal
}

func (td TDate) IsLessOrEqual(ptDate *TDate) bool {
	return td.nDateVal <= ptDate.nDateVal
}

func (td TDate) IsGreat(ptDate *TDate) bool {
	return td.nDateVal > ptDate.nDateVal
}

func (td TDate) IsGreatOrEqual(ptDate *TDate) bool {
	return td.nDateVal >= ptDate.nDateVal
}

type TTime struct {
	nTimeVal int32
}

func (tt TTime) IsValid() bool {
	return tt.nTimeVal < MsecsPerDay
}

func (tt TTime) isValid(nHour, nMin, nSec, nMs int32) bool {
	return (nHour < 24 && nMin < 60 && nSec < 60 && nMs < 1000)
}

// 将当前时间设置为指定时间
func (tt *TTime) SetTime(nHour, nMin, nSec, nMs int32) bool {
	if !tt.isValid(nHour, nMin, nSec, nMs) {
		return false
	} else {
		tt.nTimeVal = (nHour*SecsPerHour+nMin*SecsPerMin+nSec)*MsecsPerSec + nMs
		return true
	}
}

// 返回当前时间
func CurTime() TTime {
	tNow := time.Now()
	nHour, nMinute, nSecond := tNow.Clock()
	return TTime{int32(nHour)*MsecsPerHour + int32(nMinute)*MsecsPerMin + int32(nSecond)*MsecsPerSec + int32(tNow.Nanosecond())/1000000}
}

// 生成默认时间，即零点时间
func NewDefaultTime() *TTime {
	return &TTime{0}
}

// 生成指定时间
func NewSpecTime(nHour, nMin, nSec, nMs int32) (*TTime, error) {
	tRetTime := &TTime{0}

	if !tRetTime.SetTime(nHour, nMin, nSec, nMs) {
		tRetTime = nil
		return nil, errors.New("invalid param")
	}

	return tRetTime, nil
}

// 按照%02d:%02d:%02d格式格式创建字符串指定时间(精确到秒)
func NewNormalTime(strTime string) (*TTime, error) {
	var nHour, nMin, nSec int32

	if _, err := fmt.Sscanf(strTime, "%02d:%02d:%02d", &nHour, &nMin, &nSec); err != nil {
		return nil, err
	}

	tRetTime := &TTime{0}
	if !tRetTime.SetTime(nHour, nMin, nSec, 0) {
		tRetTime = nil
		return nil, errors.New("invalid time")
	}

	return tRetTime, nil
}

// 按照%02d:%02d:%02d:%03d格式格式创建字符串指定时间(精确到毫秒)
func NewDetailTime(strTime string) (*TTime, error) {
	var nHour, nMin, nSec, nMs int32

	_, err := fmt.Sscanf(strTime, "%02d:%02d:%02d:%03d", &nHour, &nMin, &nSec, &nMs)
	if err != nil {
		return nil, err
	}

	tRetTime := &TTime{0}
	if !tRetTime.SetTime(nHour, nMin, nSec, nMs) {
		tRetTime = nil
		return nil, errors.New("invalid time")
	}

	return tRetTime, nil
}

// 返回当前时间中的小时数
func (tt TTime) Hour() int32 {
	return tt.nTimeVal / MsecsPerHour
}

// 返回当前时间中的分钟数
func (tt TTime) Minute() int32 {
	return (tt.nTimeVal % MsecsPerHour) / MsecsPerMin
}

// 返回当前时间中的秒数
func (tt TTime) Second() int32 {
	return (tt.nTimeVal / MsecsPerSec) % SecsPerMin
}

// 返回当前时间中的毫秒数
func (tt TTime) MilliSec() int32 {
	return tt.nTimeVal % MsecsPerSec
}

func (tt TTime) AddSecs(nSecs int32) TTime {
	return tt.AddMilliSecs(nSecs * MsecsPerSec)
}

func (tt TTime) AddMilliSecs(nMs int32) TTime {
	tTmpTime := TTime{0}

	if nMs < 0 {
		nNegDays := (MsecsPerDay - nMs) / MsecsPerDay
		tTmpTime.nTimeVal = (tt.nTimeVal + nMs + nNegDays*MsecsPerDay) % MsecsPerDay
	} else {
		tTmpTime.nTimeVal = (tt.nTimeVal + nMs) % MsecsPerDay
	}

	return tTmpTime
}

func (tt TTime) SecsTo(ptDstTime *TTime) int32 {
	return (ptDstTime.nTimeVal - tt.nTimeVal) / MsecsPerSec
}

func (tt TTime) MilliSecsTo(ptDstTime *TTime) int32 {
	return ptDstTime.nTimeVal - tt.nTimeVal
}

// 将时间转换为精确到毫秒的字符串，格式为%02d:%02d:%02d:%03d
func (tt TTime) ToDetailTime() string {
	return fmt.Sprintf("%02d:%02d:%02d:%03d", tt.Hour(), tt.Minute(), tt.Second(), tt.MilliSec())
}

// 将时间转换为精确到秒的字符串，格式为%02d:%02d:%02d
func (tt TTime) ToNormalTime() string {
	return fmt.Sprintf("%02d:%02d:%02d", tt.Hour(), tt.Minute(), tt.Second())
}

// 以下为时间的一些比较函数
func (tt TTime) IsEqual(ptTime *TTime) bool {
	return tt.nTimeVal == ptTime.nTimeVal
}

func (tt TTime) IsNotEqual(ptTime *TTime) bool {
	return tt.nTimeVal != ptTime.nTimeVal
}

func (tt TTime) IsLess(ptTime *TTime) bool {
	return tt.nTimeVal < ptTime.nTimeVal
}

func (tt TTime) IsLessOrEqual(ptTime *TTime) bool {
	return tt.nTimeVal <= ptTime.nTimeVal
}

func (tt TTime) IsGreat(ptTime *TTime) bool {
	return tt.nTimeVal > ptTime.nTimeVal
}

func (tt TTime) IsGreatOrEqual(ptTime *TTime) bool {
	return tt.nTimeVal >= ptTime.nTimeVal
}

type TDateTime struct {
	tDate TDate
	tTime TTime
}

func CurDateTime() TDateTime {
	tNow := time.Now()
	nYear, nMonth, nDay := tNow.Date()
	val := gregToJulian(int32(nYear), int32(nMonth), int32(nDay))
	tDate := TDate{val}

	nHour, nMinute, nSecond := tNow.Clock()
	tTime := TTime{int32(nHour)*MsecsPerHour + int32(nMinute)*MsecsPerMin + int32(nSecond)*MsecsPerSec + int32(tNow.Nanosecond())/1000000}

	return TDateTime{tDate, tTime}
}

// 生成指定数值的日期时间
func NewDateTime(nYear, nMonth, nDay, nHour, nMin, nSec, nMs int32) (*TDateTime, error) {
	ptDate, errDate := NewSpecDate(nYear, nMonth, nDay)
	if errDate != nil {
		return nil, errDate
	}

	ptTime, errTime := NewSpecTime(nHour, nMin, nSec, nMs)
	if errTime != nil {
		ptDate = nil
		return nil, errTime
	}

	return &TDateTime{*ptDate, *ptTime}, nil
}

// 根据TDate TTime生成日期时间
func NewSpecDateTime(ptDate *TDate, ptTime *TTime) *TDateTime {
	return &TDateTime{*ptDate, *ptTime}
}

// 根据%04d-%02d-%02d %02d:%02d:%02d格式的字符串生成日期时间(精确到秒)
func NewNormalDateTime(strDateTime string) (*TDateTime, error) {
	var nYear, nMonth, nDay, nHour, nMin, nSec int32

	if _, err := fmt.Sscanf(strDateTime, "%04d-%02d-%02d %02d:%02d:%02d", &nYear, &nMonth, &nDay, &nHour, &nMin, &nSec); err != nil {
		return nil, err
	}

	return NewDateTime(nYear, nMonth, nDay, nHour, nMin, nSec, 0)
}

// 根据%04d-%02d-%02d %02d:%02d:%02d:%03d格式的字符串生成日期时间(精确到毫秒)
func NewDetailDateTime(strDateTime string) (*TDateTime, error) {
	var nYear, nMonth, nDay, nHour, nMin, nSec, nMs int32

	if _, err := fmt.Sscanf(strDateTime, "%04d-%02d-%02d %02d:%02d:%02d:%03d", &nYear, &nMonth, &nDay, &nHour, &nMin, &nSec, &nMs); err != nil {
		return nil, err
	}

	return NewDateTime(nYear, nMonth, nDay, nHour, nMin, nSec, nMs)
}

// 将目前日期时间设置为指定数值
func (t *TDateTime) SetDateTime(nYear, nMonth, nDay, nHour, nMin, nSec, nMs int32) bool {
	ptDate := &t.tDate
	ptTime := &t.tTime

	return ptDate.SetDate(nYear, nMonth, nDay) && ptTime.SetTime(nHour, nMin, nSec, nMs)
}

// 返回目前日期时间加上指定天数后的日期时间，参数可正可负
func (t TDateTime) AddDays(nDays int32) TDateTime {
	return TDateTime{t.tDate.AddDays(nDays), t.tTime}
}

// 返回目前日期时间加上指定时数后的日期时间，参数可正可负
func (t TDateTime) AddHours(sdwHours int64) TDateTime {
	return t.AddMilliSecs(sdwHours * int64(MsecsPerHour))
}

// 返回目前日期时间加上指定分钟数后的日期时间，参数可正可负
func (t TDateTime) AddMinutes(sdwMins int64) TDateTime {
	return t.AddMilliSecs(sdwMins * int64(MsecsPerMin))
}

// 返回目前日期时间加上指定秒数后的日期时间，参数可正可负
func (t TDateTime) AddSecs(sdwSecs int64) TDateTime {
	return t.AddMilliSecs(sdwSecs * int64(MsecsPerSec))
}

// 返回目前日期时间加上指定毫秒数后的日期时间，参数可正可负
func (t TDateTime) AddMilliSecs(sdwMs int64) TDateTime {
	nDateVal := t.tDate.nDateVal
	nTimeVal := t.tTime.nTimeVal
	var nSign int32 = 1

	if sdwMs < 0 {
		sdwMs = -sdwMs
		nSign = -1
	}

	if sdwMs >= int64(MsecsPerDay) {
		nDateVal += nSign * int32(sdwMs/int64(MsecsPerDay))
		sdwMs %= int64(MsecsPerDay)
	}

	nTimeVal += nSign * int32(sdwMs)
	if nTimeVal < 0 {
		nTimeVal = MsecsPerDay - nTimeVal - 1
		nDateVal -= nTimeVal / MsecsPerDay
		nTimeVal %= MsecsPerDay
		nTimeVal = MsecsPerDay - nTimeVal - 1
	} else if nTimeVal >= MsecsPerDay {
		nDateVal += nTimeVal / MsecsPerDay
		nTimeVal = nTimeVal % MsecsPerDay
	}

	return TDateTime{tDate: TDate{nDateVal}, tTime: TTime{nTimeVal}}
}

// 获取目前日期时间对应的格林尼治时间
func (t TDateTime) GetUnix() uint32 {
	tSysTime := time.Date(int(t.tDate.Year()), time.Month(t.tDate.Month()), int(t.tDate.Day()),
		int(t.tTime.Hour()), int(t.tTime.Minute()), int(t.tTime.Second()), 0, time.Local)
	return uint32(tSysTime.Unix())
}

// 将目前日期时间设置为指定格林尼治时间对应的日期时间
func (t *TDateTime) SetUnix(dwUnix int64) bool {
	tSysTime := time.Unix(dwUnix, 0)
	nYear, nMonth, nDay := tSysTime.Date()
	nHour, nMinute, nSecond := tSysTime.Clock()
	return t.SetDateTime(int32(nYear), int32(nMonth), int32(nDay),
		int32(nHour), int32(nMinute), int32(nSecond), 0)
}

// 计算两个日期时间之间的相隔天数
func (t TDateTime) DaysTo(ptDateTime *TDateTime) int32 {
	return t.tDate.DaysTo(&ptDateTime.tDate)
}

// 计算两个日期时间之间的相隔秒数
func (t TDateTime) SecsTo(ptDateTime *TDateTime) int64 {
	return int64(t.tTime.SecsTo(&ptDateTime.tTime)) + int64(t.tDate.DaysTo(&ptDateTime.tDate))*int64(SecsPerDay)
}

// 计算两个日期时间之间的相隔毫秒数
func (t TDateTime) MilliSecsTo(ptDateTime *TDateTime) int64 {
	return int64(t.tTime.MilliSecsTo(&ptDateTime.tTime)) + int64(t.tDate.DaysTo(&ptDateTime.tDate))*int64(MsecsPerDay)
}

// 将目前日期时间转换为精确到秒的字符串，格式为%04d-%02d-%02d %02d:%02d:%02d
func (t TDateTime) ToNormalDateTime() string {
	return t.tDate.ToString() + " " + t.tTime.ToNormalTime()
}

// 将目前日期时间转换为精确到毫秒的字符串，格式为%04d-%02d-%02d %02d:%02d:%02d:%03d
func (t TDateTime) ToDetailDateTime() string {
	return t.tDate.ToString() + " " + t.tTime.ToDetailTime()
}

// 以下为日期时间的一些比较函数
func (t TDateTime) IsEqual(ptDateTime *TDateTime) bool {
	return t.tDate.IsEqual(&ptDateTime.tDate) && t.tTime.IsEqual(&ptDateTime.tTime)
}

func (t TDateTime) IsNotEqual(ptDateTime *TDateTime) bool {
	return t.tDate.IsNotEqual(&ptDateTime.tDate) || t.tTime.IsNotEqual(&ptDateTime.tTime)
}

func (t TDateTime) IsLess(ptDateTime *TDateTime) bool {
	if t.tDate.IsLess(&ptDateTime.tDate) {
		return true
	}

	return t.tDate.IsEqual(&ptDateTime.tDate) && t.tTime.IsLess(&ptDateTime.tTime)
}

func (t TDateTime) IsLessOrEqual(ptDateTime *TDateTime) bool {
	if t.tDate.IsLess(&ptDateTime.tDate) {
		return true
	}

	return t.tDate.IsEqual(&ptDateTime.tDate) && t.tTime.IsLessOrEqual(&ptDateTime.tTime)
}

func (t TDateTime) IsGreat(ptDateTime *TDateTime) bool {
	if t.tDate.IsGreat(&ptDateTime.tDate) {
		return true
	}

	return t.tDate.IsEqual(&ptDateTime.tDate) && t.tTime.IsGreat(&ptDateTime.tTime)
}

func (t TDateTime) IsGreatOrEqual(ptDateTime *TDateTime) bool {
	if t.tDate.IsGreat(&ptDateTime.tDate) {
		return true
	}

	return t.tDate.IsEqual(&ptDateTime.tDate) && t.tTime.IsGreatOrEqual(&ptDateTime.tTime)
}

func (t TDateTime) IsValid() bool {
	return t.tDate.IsValid() && t.tTime.IsValid()
}

func (t *TDateTime) GetDate() *TDate {
	return &t.tDate
}

func (t *TDateTime) GetTime() *TTime {
	return &t.tTime
}

func (t *TDateTime) SetDate(ptDate *TDate) {
	t.tDate = *ptDate
}

func (t *TDateTime) SetTime(ptTime *TTime) {
	t.tTime = *ptTime
}

func (t TDateTime) Year() int32 {
	return t.tDate.Year()
}

func (t TDateTime) Month() int32 {
	return t.tDate.Month()
}

func (t TDateTime) Day() int32 {
	return t.tDate.Day()
}

func (t TDateTime) Hour() int32 {
	return t.tTime.Hour()
}

func (t TDateTime) Minute() int32 {
	return t.tTime.Minute()
}

func (t TDateTime) Second() int32 {
	return t.tTime.Second()
}

func (t TDateTime) MilliSec() int32 {
	return t.tTime.MilliSec()
}

func (t TDateTime) DayZeroTime() uint32 {
	return t.GetUnix() - uint32(t.Hour()*SecsPerHour) - uint32(t.Minute()*SecsPerMin) - uint32(t.Second())
}

func WithinTime(start uint32, sep_time uint32) bool {
	sopen := time.Unix(int64(start), 0)
	sopen_hour, sopen_minute, sopen_second := sopen.Hour(), sopen.Minute(), sopen.Second()

	start_time := uint32(int(start) - sopen_hour*3600 - sopen_minute*60 - sopen_second)
	end_time := start_time + (sep_time * 24 * 3600)
	now := uint32(time.Now().Unix())

	if start <= now && now <= end_time {
		return true
	}
	return false
}

// 获取当天的零点时间
func TheDayZeroTime(start int64) int64 {
	var todayZeroTime int64
	sopen := time.Unix(start, 0)
	sopen_hour, sopen_minute, sopen_second := int64(sopen.Hour()), int64(sopen.Minute()), int64(sopen.Second())
	todayZeroTime = start - sopen_hour*3600 - sopen_minute*60 - sopen_second
	return todayZeroTime
}

func DaysToFirstOpenTime(openTime int64) uint32 {
	ptOpenDate := &TDateTime{}
	ptOpenDate.SetUnix(openTime)

	//获取当前时间的Date格式
	nowDate := CurDate()

	//计算间隔时间（负数）
	intervalDay := nowDate.DaysTo(ptOpenDate.GetDate())
	return uint32(math.Abs(float64(intervalDay)))
}

func IsSameDayToToday(otherTime int64) bool {
	otherDate := &TDateTime{}
	otherDate.SetUnix(otherTime)

	//当前时间的Date格式
	nowDate := CurDate()

	result := nowDate.IsEqual(otherDate.GetDate())
	return result
}

// 计算两个时间戳的相隔天数
func DaysOfTwoDays(timeone, timetwo int64) int32 {
	ptOneDate := &TDateTime{}
	ptOneDate.SetUnix(timeone)

	ptTwoDate := &TDateTime{}
	ptTwoDate.SetUnix(timetwo)

	return ptOneDate.GetDate().DaysTo(ptTwoDate.GetDate())
}
