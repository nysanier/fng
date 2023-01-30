package pkgenv

import (
	"os"

	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

func getEnvWithDefault(key, defVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defVal
}

// 加载env
func LoadEnv() {
	env := getEnvWithDefault("fn_env", pkgvar.FnEnv_Dev)
	pkgvar.FnEnv = env

	aesKeyStr := os.Getenv("fn_aes_key") // env中自行补齐到16位
	pkgvar.FnAesKey = []byte(aesKeyStr)
	if aesKeyStr == "" {
		panic(aesKeyStr)
	}

	skBuf, err := pkgfunc.DecryptByAES(pkgvar.OriSK)
	if err != nil {
		panic(err)
	}

	sk := string(skBuf)
	pkgvar.SetSK(sk)

	akBuf, err := pkgfunc.DecryptByAES(pkgvar.OriAK)
	if err != nil {
		panic(err)
	}

	ak := string(akBuf)
	pkgvar.SetAK(ak)
}
