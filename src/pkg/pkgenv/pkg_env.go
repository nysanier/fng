package pkgenv

import (
	"os"
)

// 通过操作系统env方式导入
var (
	fnEnv    string // fn_env
	fnAesKey []byte // fn_aes_key
)

// FnEnv常量定义
const (
	FnEnv_Dev   = "dev"
	FnEnv_Daily = "daily"
	FnEnv_Stg   = "stg"
	FnEnv_Prod  = "prod"
)

func GetAesKey() []byte {
	return fnAesKey
}

func GetEnv() string {
	return fnEnv
}

func IsDevEnv() bool {
	return fnEnv == FnEnv_Dev
}

func getEnvWithDefault(key, defVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defVal
}

// 加载env
func InitEnv() {
	fnEnv = getEnvWithDefault("fn_env", FnEnv_Dev)

	aesKeyStr := os.Getenv("fn_aes_key") // env中自行补齐到16位
	if aesKeyStr == "" {
		panic(aesKeyStr)
	}

	fnAesKey = []byte(aesKeyStr)
}
