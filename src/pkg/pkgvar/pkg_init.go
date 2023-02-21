package pkgvar

import "time"

// 程序启动时需要初始化的变量

// 程序启动时间
var (
	FnStartTime = "program start time"
)

// 通过操作系统env方式导入
var (
	FnEnv    string // fn_env
	FnAesKey []byte // fn_aes_key
)

// 时区信息
var (
	TzLoc *time.Location
)
