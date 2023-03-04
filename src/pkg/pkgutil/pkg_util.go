package pkgutil

import "github.com/nysanier/fng/src/pkg/pkgfunc"

var (
	fnStartTime = "program start time"
)

func getCRFC3339CstTimeStr() string {
	t := GetCstNow()
	str := pkgfunc.GetRFC3339TimeStr(t)
	return str
}

func InitStartTime() {
	fnStartTime = getCRFC3339CstTimeStr()
}

func GetStartTime() string {
	return fnStartTime
}
