package fnlog

import (
	"testing"
)

func Test_Infof(t *testing.T) {
	SetLogLevel(LogLevelInfo)
	Infof("1 info abc%v", "def")
	SetLogLevel(LogLevelWarn)
	Infof("2 warn abc%v", "def")
	SetLogLevel(LogLevelDebug)
	Infof("3 debug abc%v", "def")

	// 可以观察到1和3, 比如:
	// 2022/04/04 08:36:55 [INFO] 1 info abcdef
	// 2022/04/04 08:36:55 [INFO] 3 debug abcdef
}

func Test_Printf(t *testing.T) {
	SetLogLevel(LogLevelInfo)
	Printf(LogLevelInfo, "1 info xxx%v", "yyy")

	// 比如
	// 2022/04/04 08:36:55 [INFO] 1 info xxxyyy
}
