package pkgfunc

import (
	"log"
	"math/rand"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgvar"
)

// docker中容器，默认使用的是GMT(+00:00), 因此最好改为CST(+08:00), 方便调试
func GetCstNow() time.Time {
	t := time.Now()
	return ToCstTime(t)
}

func ToCstTime(t time.Time) time.Time {
	loc := pkgvar.TzLoc
	t2 := t.In(loc)
	return t2
}

func LoadTzLoc() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("time.LoadLocation fail: err=%v", err)
		loc = time.UTC // 默认使用 utc
	}

	return loc
}

func InitRand() {
	rand.Seed(time.Now().UnixNano())
}
