package pkgutil

import (
	"time"

	"github.com/nysanier/fng/src/pkg/pkglog"
)

// 时区信息
var (
	tzLoc *time.Location
)

func LoadTzLoc() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		pkglog.Warnv("EvtTimeLoadLocationFail",
			"Error", err)
		loc = time.UTC // 默认使用 utc
	}

	return loc
}

func InitTzLoc() {
	tzLoc = LoadTzLoc()
}

// docker中容器，默认使用的是GMT(+00:00), 因此最好改为CST(+08:00), 方便调试
func GetCstNow() time.Time {
	t := time.Now()
	return ToCstTime(t)
}

func ToCstTime(t time.Time) time.Time {
	t2 := t.In(tzLoc)
	return t2
}
