package crontab

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CrontabMe struct representing cron table
type CrontabMe struct {
	C         chan time.Time
	jobs      []*job
	wg        sync.WaitGroup
	CloseChan chan bool
	sync.RWMutex
}

// job in cron table
type job struct {
	min       uint64
	hour      uint64
	day       uint64
	month     uint64
	dayOfWeek uint64

	fn   func(args ...interface{})
	args []interface{}
	sync.RWMutex
}

// tick is individual tick that occures each minute
type tick struct {
	min       int
	hour      int
	day       int
	month     int
	dayOfWeek int
}

// New initializes and returns new cron table
func NewMe(t int64) *CrontabMe {
	ct := &CrontabMe{
		//Ticker: time.NewTicker(t),
		jobs:      []*job{},
		C:         make(chan time.Time),
		CloseChan: make(chan bool),
	}
	ct.run(t)
	return ct
}

// new creates new crontab, arg provided for testing purpose

func (ct *CrontabMe) run(serverT int64) {
	ct.wg.Add(1)

	go func() {
		const STEP int64 = 60
		serverTime := time.Unix(serverT, 0)
		systemT := time.Now().Unix()
		diff := serverT - systemT

		d := STEP - int64(serverTime.Second()) // 服务器时间下一个整分

		timer := time.NewTimer(time.Second * time.Duration(d))

		for {
			select {
			case c := <-timer.C:
				c = c.Add(time.Second * time.Duration(diff))
				d = STEP - int64(c.Second())

				ct.C <- c

				timer.Reset(time.Second * time.Duration((d)))
			case <-ct.CloseChan:
				ct.wg.Done()
				return
			}
		}
	}()
}

func (c *CrontabMe) AddJob(schedule string, fn func(...interface{}), args ...interface{}) error {
	j, err := parseSchedule(schedule)
	c.Lock()
	defer c.Unlock()
	if err != nil {
		return err
	}

	if fn == nil || reflect.ValueOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("cron job must be func()")
	}

	fnType := reflect.TypeOf(fn)
	if len(args) != fnType.NumIn() {
		return fmt.Errorf("number of func() params and number of provided params doesn't match")
	}

	for i := 0; i < fnType.NumIn(); i++ {
		a := args[i]
		t1 := fnType.In(i)
		t2 := reflect.TypeOf(a)

		if t1 != t2 {
			if t1.Kind() != reflect.Interface {
				return fmt.Errorf("param with index %d shold be `%s` not `%s`", i, t1, t2)
			}
			if !t2.Implements(t1) {
				return fmt.Errorf("param with index %d of type `%s` doesn't implement interface `%s`", i, t2, t1)
			}
		}
	}

	// all checked, add job to cron tab
	j.fn = fn
	j.args = args
	c.jobs = append(c.jobs, j)
	return nil
}

func (c *CrontabMe) MustAddJob(schedule string, fn func(...interface{}), args ...interface{}) {
	if err := c.AddJob(schedule, fn, args...); err != nil {
		panic(err)
	}
}

func (c *CrontabMe) Clear() {
	c.Lock()
	c.jobs = []*job{}
	c.Unlock()
}

func (c *CrontabMe) Close() {
	c.CloseChan <- true
}

func (c *CrontabMe) RunAll() {
	c.RLock()
	defer c.RUnlock()
	for _, j := range c.jobs {
		go j.run()
	}
}

// RunScheduled jobs
func (c *CrontabMe) RunScheduled(t time.Time) {
	fmt.Println("=== run scheduled === t ", t.Format("2006-01-02 15:04:05"), " nnn:", time.Now().Format("2006-01-02 15:04:05"))
	tick := getTick(t)
	c.RLock()
	defer c.RUnlock()

	for _, j := range c.jobs {
		if j.tick(tick) {
			j.run()
		}
	}
}

// run the job using reflection
// Recover from panic although all functions and params are checked by AddJob, but you never know.
func (j *job) run() {
	j.RLock()
	defer func() {
		if r := recover(); r != nil {
			log.Println("CrontabMe error", r)
		}
	}()
	v := reflect.ValueOf(j.fn)
	rargs := make([]reflect.Value, len(j.args))
	for i, a := range j.args {
		rargs[i] = reflect.ValueOf(a)
	}
	j.RUnlock()
	v.Call(rargs)
}

func (j *job) tick(t tick) bool {
	j.RLock()
	defer j.RUnlock()

	if j.min&uint64(1<<t.min) == 0 {
		return false
	}

	if j.hour&uint64(1<<t.hour) == 0 {
		return false
	}

	if j.day&uint64(1<<t.day) == 0 {
		return false
	}

	if j.dayOfWeek&uint64(1<<t.dayOfWeek) == 0 {
		return false
	}

	if j.month&uint64(1<<t.month) == 0 {
		return false
	}

	return true
}

// regexps for parsing schedule string
var (
	matchSpaces = regexp.MustCompile(`\s+`)
	matchN      = regexp.MustCompile(`(.*)/(\d+)`)
	matchRange  = regexp.MustCompile(`^(\d+)-(\d+)$`)
)

// parseSchedule string and creates job struct with filled times to launch, or error if synthax is wrong
func parseSchedule(s string) (*job, error) {
	var err error
	j := &job{}
	j.Lock()
	defer j.Unlock()

	s = matchSpaces.ReplaceAllLiteralString(s, " ")
	parts := strings.Split(s, " ")
	if len(parts) != 5 {
		return j, errors.New("schedule string must have five components like * * * * *")
	}

	j.min, err = parsePart(parts[0], 0, 59)
	if err != nil {
		return j, err
	}

	j.hour, err = parsePart(parts[1], 0, 23)
	if err != nil {
		return j, err
	}

	j.day, err = parsePart(parts[2], 1, 31)
	if err != nil {
		return j, err
	}

	j.month, err = parsePart(parts[3], 1, 12)
	if err != nil {
		return j, err
	}

	j.dayOfWeek, err = parsePart(parts[4], 0, 6)
	if err != nil {
		return j, err
	}

	//  day/dayOfWeek combination
	//switch {
	//case len(j.day) < 31 && len(j.dayOfWeek) == 7: // day set, but not dayOfWeek, clear dayOfWeek
	//	j.dayOfWeek = make(map[int]struct{})
	//case len(j.dayOfWeek) < 7 && len(j.day) == 31: // dayOfWeek set, but not day, clear day
	//	j.day = make(map[int]struct{})
	//default:
	//	// both day and dayOfWeek are * or both are set, use combined
	//	// i.e. don't do anything here
	//}

	return j, nil
}

// parsePart parse individual schedule part from schedule string
func parsePart(s string, min, max int) (uint64, error) {
	var r uint64

	// wildcard pattern
	if s == "*" {
		for i := min; i <= max; i++ {
			r |= 1 << i
		}
		return r, nil
	}

	// */2 1-59/5 pattern
	if matches := matchN.FindStringSubmatch(s); matches != nil {
		localMin := min
		localMax := max
		if matches[1] != "" && matches[1] != "*" {
			if rng := matchRange.FindStringSubmatch(matches[1]); rng != nil {
				localMin, _ = strconv.Atoi(rng[1])
				localMax, _ = strconv.Atoi(rng[2])
				if localMin < min || localMax > max {
					return 0, fmt.Errorf("out of range for %s in %s. %s must be in range %d-%d", rng[1], s, rng[1], min, max)
				}
			} else {
				return 0, fmt.Errorf("unable to parse %s part in %s", matches[1], s)
			}
		}
		n, _ := strconv.Atoi(matches[2])
		for i := localMin; i <= localMax; i += n {
			r |= 1 << i
		}
		return r, nil
	}

	// 1,2,4  or 1,2,10-15,20,30-45 pattern
	parts := strings.Split(s, ",")
	for _, x := range parts {
		if rng := matchRange.FindStringSubmatch(x); rng != nil {
			localMin, _ := strconv.Atoi(rng[1])
			localMax, _ := strconv.Atoi(rng[2])
			if localMin < min || localMax > max {
				return 0, fmt.Errorf("out of range for %s in %s. %s must be in range %d-%d", x, s, x, min, max)
			}
			for i := localMin; i <= localMax; i++ {
				r |= 1 << i
			}
		} else if i, err := strconv.Atoi(x); err == nil {
			if i < min || i > max {
				return 0, fmt.Errorf("out of range for %d in %s. %d must be in range %d-%d", i, s, i, min, max)
			}
			r |= 1 << i
		} else {
			return 0, fmt.Errorf("unable to parse %s part in %s", x, s)
		}
	}

	return r, nil
}

// getTick returns the tick struct from time
func getTick(t time.Time) tick {
	return tick{
		min:       t.Minute(),
		hour:      t.Hour(),
		day:       t.Day(),
		month:     int(t.Month()),
		dayOfWeek: int(t.Weekday()),
	}
}
