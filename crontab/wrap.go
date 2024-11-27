package crontab

import (
	"sync"
	"time"

	"github.com/Lei2050/lei-utils/cls"
)

/*
类似linux的crontab功能，目前不支持"* /2"语法。
C用来做上层的串行化。
usage:

	cron := crontab.New()
	cron.Run(time.Now().Unix())
	cron.AddSchedule(...)

	for range cron.C {
		cron.Process()
	}
*/
type Crontab struct {
	cron crontab

	C  chan struct{}
	es map[uint64]struct{} //即将要执行的
	wg sync.WaitGroup

	cls.CloseUtil
}

func New() *Crontab {
	return &Crontab{
		cron: crontab{
			entries:     make(map[uint64]*Entry),
			idmgr:       100000000,
			lastRunTime: 0,
		},
		C:         make(chan struct{}, 10),
		es:        make(map[uint64]struct{}),
		CloseUtil: cls.MakeCloseUtil(),
	}
}

func (c *Crontab) Remove(id uint64) {
	c.cron.remove(id)
}

// 如果提供了id，则会尽量使用该id，如果不重复的话；
// 否则自动返回一个id。
func (c *Crontab) AddSchedule(id uint64,
	minute, hour, dom, month, dow string,
	cmd func(...interface{}), args ...interface{}) (uint64, error) {
	return c.cron.addSchedule(id, minute, hour, dom, month, dow, cmd, args...)
}

func (c *Crontab) AddScheduleByStr(id uint64, str string,
	cmd func(...interface{}), args ...interface{}) (uint64, error) {
	return c.cron.addScheduleByStr(id, str, cmd, args...)
}

func (ctb *Crontab) run(tm int64) {
	ctb.wg.Add(1)

	go func() {
		const STEP int64 = 60 //最低细粒度60秒，目前只支持到每分钟
		t := time.Unix(tm, 0)
		now := time.Now().Unix()
		diff := tm - now

		var d int64 = STEP - int64(t.Second())
		//fmt.Printf("d:%d\n", d)
		timer := time.NewTimer(time.Second * time.Duration(d))

		for {
			select {
			case c := <-timer.C:
				c = c.Add(time.Second * time.Duration(diff))
				d = STEP - int64(c.Second())
				//fmt.Printf("d:%d, c:%d\n", d, c.Unix())

				eslen := len(ctb.es)
				ctb.cron.findJobs(ctb.es, c)
				if eslen == 0 && len(ctb.es) > 0 {
					//fmt.Printf("notify %d !\n", len(ctb.es))
					ctb.C <- struct{}{}
				}

				timer.Reset(time.Second * time.Duration(d))
			case <-ctb.C:
				ctb.wg.Done()
				return
			}
		}
	}()
}

// 外部提供一个当前时间戳tm，时间判断是依据该时间戳为基准，而不是标准包的时间；
// 因为考虑到上层应用可能会有自己的时间系统（不是按标准包的时间）。
// 调用Run会先停止当前的Crontab，然后再重新启动。若之前有未完成的任务，则会丢弃。
func (c *Crontab) Run(tm int64) {
	//先关闭原来的，再开新的一个。
	//考虑到上层修改时间重新Run
	if !c.IsClosed() {
		c.Close(nil)
		c.wg.Wait() //等待上一个run的关闭
	}
	c.CloseUtil = cls.MakeCloseUtil()

	c.es = make(map[uint64]struct{})
	close(c.C)
	c.C = make(chan struct{}, 10)

	c.run(tm)
}

func (c *Crontab) Process() {
	for id := range c.es {
		c.cron.do(id)
	}
	c.es = make(map[uint64]struct{})
}
