package pkgfunc

import "time"

func GetRFC3339TimeStr(t time.Time) string {
	str := t.Format(time.RFC3339)
	return str
}
