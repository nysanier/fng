package pkgvar

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

func GetRegionID() string {
	return RegionID
}
