package crontab

import (
	"fmt"
	"strings"
	"time"
)

type crontab struct {
	entries map[uint64]*Entry
	idmgr   uint64

	lastRunTime int64
}

func (c *crontab) remove(id uint64) {
	delete(c.entries, id)
}

func (c *crontab) genid(id uint64) uint64 {
	if id <= 0 {
		id = c.idmgr + 1
	}
	for _, exists := c.entries[id]; exists; _, exists = c.entries[id] {
		id++
	}
	c.idmgr = id
	return id
}

func (c *crontab) addSchedule(id uint64,
	minute, hour, dom, month, dow string,
	cmd func(...interface{}), args ...interface{}) (uint64, error) {
	//c.idmgr++
	id = c.genid(id)

	e, err := newEntry(id, minute, hour, dom, month, dow, cmd, args...)
	if err != nil {
		return 0, err
	}
	if e == nil {
		return 0, fmt.Errorf("alloc entry failed")
	}

	c.entries[id] = e
	return id, nil
}

func (c *crontab) addScheduleByStr(id uint64, str string,
	cmd func(...interface{}), args ...interface{}) (uint64, error) {

	strs := strings.Split(str, " ")
	if len(strs) != 5 {
		return 0, fmt.Errorf("error input:%s", str)
	}

	return c.addSchedule(id, strs[0], strs[1], strs[2], strs[3], strs[4], cmd, args...)
}

func (c *crontab) findJobs(es map[uint64]struct{}, t time.Time) {
	for _, e := range c.entries {
		if e.Test(t) {
			es[e.id] = struct{}{}
		}
	}
}

func (c *crontab) do(id uint64) {
	e := c.entries[id]
	if e == nil {
		return
	}
	if e.cmd != nil {
		e.cmd(e.args...)
	}
}
