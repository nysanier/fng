package pkgvar

// FnEnv常量定义
const (
	FnEnv_Dev   = "dev"
	FnEnv_Daily = "daily"
	FnEnv_Stg   = "stg"
	FnEnv_Prod  = "prod"
)

func IsDevEnv() bool {
	return FnEnv == FnEnv_Dev
}

const (
	// TODO: 通过ots/kms等存储和加解密， aksk需要env常量来解码
	OriAK    = "QvBLI9DsD+kwwFnK9xlBqO/kzAwDArMGV3TFnyxUxkY="
	OriSK    = "zumURn5+uAyoSMbSvHbDp46fFI9x+jIgpfQ0ks7C8CM="
	RegionID = "cn-hangzhou"
)

var (
	ak_ = ""
	sk_ = ""
)

// 变量需要通过函数来获取
func GetAK() string {
	return ak_
}

func SetAK(ak string) {
	ak_ = ak
}

func GetSK() string {
	return sk_
}

func SetSK(sk string) {
	sk_ = sk
}

// 配置表有一个value列
const (
	ConfigFieldPk1   = "block"
	ConfigFieldPk2   = "section"
	ConfigFieldValue = "value"
)
