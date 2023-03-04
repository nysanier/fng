package pkglog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/nysanier/fng/src/pkg/pkgvar"
)

type LogIntf interface {
	WriteLog(keys []string, vals []string)
}

var logIntf LogIntf

func SetImpl(intf LogIntf) {
	logIntf = intf
}

func InitLog() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func shortFilePath(fp string) string {
	l := strings.Split(fp, "/")
	ln := len(l)
	switch ln {
	case 0:
		return ""
	case 1:
		return l[0]
	case 2:
		return strings.Join([]string{l[0], l[1]}, "/")
	default:
		return strings.Join([]string{l[ln-3], l[ln-2], l[ln-1]}, "/")
	}
}

func writeLog(level, event string, items ...interface{}) {
	var itemList []interface{}
	itemList = append(itemList, items...)
	if len(items)%2 == 1 {
		itemList = append(itemList, "") // 填补空串
	}

	writeToLocal(level, event, itemList)

	//timeStr := time.Now().Format(time.RFC3339)

	// 调用栈信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// 元信息
	itemList = append(itemList,
		//"_fn_rfc3339", timeStr,
		"_fn_env", pkgvar.FnEnv,
		"_fn_level", level,
		"_fn_file", shortFilePath(file),
		"_fn_line", line,
		"Event", event, // TODO: 通过开关来定义？
		// TODO: 通过开关带上RequestID信息？
	)

	var keys, vals []string
	for i := 0; i < len(itemList); i += 2 {
		k := intfToString(itemList[i])
		v := intfToString(itemList[i+1])
		keys = append(keys, k)
		vals = append(vals, v)
	}

	logIntf.WriteLog(keys, vals)
}

func writeToLocal(level, event string, itemList []interface{}) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("[%v] %v", level, event))
	for i := 0; i < len(itemList); i += 2 {
		k := intfToString(itemList[i])
		v := intfToString(itemList[i+1])
		buf.WriteString(", " + k + "=" + v)
	}

	log.Default().Output(4, fmt.Sprintf("%v", buf.String()))
}

func intfToString(v interface{}) string {
	switch v.(type) {
	case int, uint, uintptr, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return fmt.Sprintf("%v", v)
	case string, complex64, complex128:
		return fmt.Sprintf("%v", v)
	//case byte: // 和uint8重复
	//	return fmt.Sprintf("%v", v)
	//case rune: // 和int32重复
	//	return fmt.Sprintf("%v", v)
	default:
		buf, _ := json.Marshal(v)
		return string(buf)
	}
}

// TODO: 根据开关，决定哪些日志不用输出
func Debugv(event string, items ...interface{}) {
	writeLog("DEBUG", event, items...)
}

func Infov(event string, items ...interface{}) {
	writeLog("INFO", event, items...)
}

func Warnv(event string, items ...interface{}) {
	writeLog("WARN", event, items...)
}

func Errorv(event string, items ...interface{}) {
	writeLog("ERROR", event, items...)
}

func Fatalv(event string, items ...interface{}) {
	writeLog("FATAL", event, items...)
	panic("fatal log")
}
