package models

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"sync/atomic"
	"time"
)

const (
	// TimeTypeDay 这里的时间格式为  05/04/15/02/01  表示 5 月 4 日的 15:02:01。 Day 只能做一次
	TimeTypeDay = 1
	// TimeTypeWeek 这里的时间格式为 15:02:01 表示15:02:01 args中的值 从 1~7 可以做多次 若多次 则表示每月的周几去做
	TimeTypeWeek  = 2
	TimeTypeMonth = 3
	TimeTypeYear  = 4

	OneDaySecond = 24 * 3600

	StrSplitChar = "/"
)

var (
	// 秒/分钟/小时/天 of month / 月 / day of week
	cronTemplate = []string{"*", "*", "*", "*", "*", "?"}
)

// Job ===================================================

type DefaultJob struct {
	dbID      int64
	cronID    int32
	Handler   func()
	ExecCount int64
	EndTime   int64
}

func NewNormalJob(f func(), endTime int64) *DefaultJob {
	return &DefaultJob{
		EndTime: endTime,
		cronID:  0,
		Handler: f,
	}
}

func (nj *DefaultJob) SetID(id int32) {
	nj.cronID = id
}

func (nj *DefaultJob) Run() {
	nowTime := time.Now().Unix()

	if nj.EndTime != 0 && nowTime > nj.EndTime {
		err := RemoveJob(int(nj.cronID))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	atomic.AddInt64(&(nj.ExecCount), 1)
	nj.Handler()
}

// cron ===========================================

type CronTimer struct {
	c *cron.Cron
}

var (
	cronTimer *CronTimer
)

func InitCronTimer() {
	cronTimer = newCronTimer()
}

func newCronTimer() *CronTimer {
	ct := CronTimer{}
	ct.c = cron.New(cron.WithSeconds())
	ct.c.Start()

	return &ct
}

func AddCron(cronStr string, endTime int64, f func()) error {

	if cronTimer == nil {
		return errors.New("please init cron timer")
	}

	job := NewNormalJob(f, endTime)
	cronID, err := cronTimer.c.AddJob(cronStr, job)

	if err != nil {
		return err
	}
	job.SetID(int32(cronID))
	return nil
}

func RemoveJob(jobID int) error {
	cronTimer.c.Remove(cron.EntryID(jobID))
	fmt.Println("cron remove job: cronID is ", jobID)
	return nil
}
