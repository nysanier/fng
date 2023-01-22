package pkgfunc

import (
	"log"
	"time"
)

// docker中容器，默认使用的是GMT(+00:00), 因此最好改为CST(+08:00), 方便调试
func GetCstNow() time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("time.LoadLocation fail: err=%v", err)
		loc = time.UTC
	}

	t := time.Now().In(loc)
	return t
}
