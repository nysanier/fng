package pkgfunc

import "encoding/json"

func FormatJson(v interface{}) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}
