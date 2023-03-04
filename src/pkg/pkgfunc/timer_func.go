package pkgfunc

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nysanier/fng/src/pkg/pkglog"
)

type TimerFunc func() error

// 基于两次任务之间间隔的定时器
type Timer struct {
	TickerCheck *time.Ticker // 内部每秒都检查一次是否要求stop
	Started     int32
	Wg          sync.WaitGroup
	Interval    time.Duration // 两次执行的间隔时间
	Func        TimerFunc
	FirstDelay  time.Duration // 首次执行delay时间，默认会立刻执行一次
}

func NewTimer(f TimerFunc, interval time.Duration) *Timer {
	p := &Timer{
		Interval: interval,
		Func:     f,
	}

	return p
}

func (p *Timer) SetFirstDelay(firstDelay time.Duration) *Timer {
	p.FirstDelay = firstDelay
	return p
}

// 启动一个协程执行定时器任务
func (p *Timer) Start() {
	p.Started = 1
	p.Wg.Add(1)
	go p.Run()
}

// 等待结束
func (p *Timer) Stop() {
	atomic.StoreInt32(&p.Started, 0)
	p.Wg.Wait()
}

func (p *Timer) Run() {
	if int64(p.FirstDelay) > 0 {
		time.Sleep(p.FirstDelay) // 首次执行延迟
	}

	// 启动的时候先执行一次
	if err := p.Func(); err != nil {
		pkglog.Warnv("EvtTimerRunFuncFirstlyFail",
			"Error", err)
		return
	}

	for {
		started := atomic.LoadInt32(&p.Started)
		if started == 0 {
			break
		}

		time.Sleep(p.Interval)
		if err := p.DoTimerFunc(); err != nil {
			pkglog.Warnv("EvtTimerRunFuncFail",
				"Error", err)
			continue
		}
	}

	p.Wg.Done()
}

func (p *Timer) DoTimerFunc() (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err3 := fmt.Errorf("DoTimerFunc panic")
			pkglog.Warnv("EvtTimerRunFuncPanic",
				"error", err2) // 真实原因只在这里打印
			err = err3 // 将panic信息传递出去
		}
	}()

	return p.Func()
}
