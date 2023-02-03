package pkgfunc

import (
	"log"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Timer(t *testing.T) {
	patches := NewPatches()
	f := func() error {
		log.Printf("x")
		return nil
	}
	timer := NewTimer(f, time.Second).SetFirstDelay(time.Second * 2)
	Convey("Test_Timer", t, func() {
		defer patches.Reset()

		Convey("ok", func() {
			defer patches.Reset()

			log.Printf("starting timer ...")
			timer.Start()
			log.Printf("start timer ok")
			time.Sleep(time.Second * 8)

			log.Printf("stoping timer ...")
			timer.Stop()
			log.Printf("stop timer ok")
		})
	})
}

/*
2023/02/03 22:49:12 starting timer ...
2023/02/03 22:49:12 start timer ok
2023/02/03 22:49:14 x (首次执行等待了2秒）
2023/02/03 22:49:15 x
2023/02/03 22:49:16 x
2023/02/03 22:49:17 x
2023/02/03 22:49:18 x
2023/02/03 22:49:19 x
2023/02/03 22:49:20 stoping timer ...
2023/02/03 22:49:20 x (请求stop的时候，还执行了一次)
2023/02/03 22:49:20 stop timer ok

*/
