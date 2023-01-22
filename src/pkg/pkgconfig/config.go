package pkgconfig

import (
	"os"

	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

var (
	Env = Env_Dev
)

const (
	Env_Dev   = "dev"
	Env_Daily = "daily"
	Env_Stg   = "stg"
	Env_Prod  = "prod"
)

var (
	// TODO: 通过ots/kms等存储和加解密
	AK = "QvBLI9DsD+kwwFnK9xlBqO/kzAwDArMGV3TFnyxUxkY="
	SK = "zumURn5+uAyoSMbSvHbDp46fFI9x+jIgpfQ0ks7C8CM="
)

func LoadConfig() {
	Env = os.Getenv("fn_env")
	if Env == "" {
		Env = Env_Dev
	}

	aesKey := os.Getenv("fn_aes_key") // env中自行补齐到16位
	pkgvar.AesKey = []byte(aesKey)

	if aesKey == "" {
		panic(aesKey)
	}

	skBuf, err := pkgfunc.DecryptByAES(SK)
	if err != nil {
		panic(err)
	}

	SK = string(skBuf)

	akBuf, err := pkgfunc.DecryptByAES(AK)
	if err != nil {
		panic(err)
	}

	AK = string(akBuf)

}

//func GetConfig(block, section, item string) interface{} {
//
//}

//func GetConfigString(block, section, item string) string {
//	v := GetConfig(block, section, item)
//	ret := v.(string)
//	return ret
//}
