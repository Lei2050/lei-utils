package crontab

import (
	"testing"
	"time"
)

func TestCrontab(t *testing.T) {
	cron := New()
	cron.AddSchedule(0, "1", "12", "", "", "", func(args ...interface{}) {
		t.Logf("%s 1 12 * * *", args[0])
	}, "hello")
	cron.AddSchedule(0, "", "", "", "", "", func(args ...interface{}) {
		t.Logf("1 * * * *")
	})
	cron.AddSchedule(0, "0", "5", "", "", "", func(args ...interface{}) {
		t.Logf("* 5 * * *")
	})
	cron.AddSchedule(0, "0", "5", "", "", "2-4", func(args ...interface{}) {
		t.Logf("* 5 * * 2-4")
	})
	cron.AddSchedule(0, "0", "5", "10", "1", "", func(args ...interface{}) {
		t.Logf("* 5 10 1 *")
	})
	cron.AddSchedule(0, "0", "7", "4", "", "", func(args ...interface{}) {
		t.Logf("* 7 4 * *")
	})
	cron.AddSchedule(0, "0", "7", "4", "", "", func(args ...interface{}) {
		t.Logf("* 7 4 * *")
	})
	cron.AddSchedule(0, "30", "6,12", "", "", "", func(args ...interface{}) {
		t.Logf("30 6,12 * * *")
	})

	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Now())
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 30, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 30, 5, 30, 1, 0, time.UTC))
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 30, 12, 1, 1, 0, time.UTC))
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2021, 1, 1, 1, 1, 0, 0, time.UTC))
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 24, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 25, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 26, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2021, 8, 27, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2020, 1, 10, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2020, 1, 11, 5, 0, 0, 0, time.UTC))
	cron.Process()
	t.Logf("================================")
	cron.cron.findJobs(cron.es, time.Date(2020, 1, 11, 6, 30, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2020, 1, 11, 12, 30, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2020, 1, 11, 6, 31, 0, 0, time.UTC))
	cron.Process()
	t.Logf("===")
	cron.cron.findJobs(cron.es, time.Date(2020, 1, 11, 7, 30, 0, 0, time.UTC))
	cron.Process()
	t.Logf("================================")

	cron.cron.findJobs(cron.es, time.Unix(1630313820, 0))
	cron.Process()
	cron.cron.findJobs(cron.es, time.Unix(1630314060, 0))
	cron.Process()

	cron.Run(time.Now().Unix())

	for range cron.C {
		t.Logf("%d", time.Now().Unix())
		cron.Process()
	}
}
