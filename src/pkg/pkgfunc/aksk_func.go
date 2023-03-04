package pkgfunc

import (
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

const (
	// TODO: 通过ots/kms等存储和加解密， aksk需要env常量来解码
	oriAK = "QvBLI9DsD+kwwFnK9xlBqO/kzAwDArMGV3TFnyxUxkY="
	oriSK = "zumURn5+uAyoSMbSvHbDp46fFI9x+jIgpfQ0ks7C8CM="
)

func InitAksk() {
	skBuf, err := DecryptByAES(oriSK)
	if err != nil {
		panic(err)
	}

	pkgvar.SetSK(string(skBuf))

	akBuf, err := DecryptByAES(oriAK)
	if err != nil {
		panic(err)
	}

	pkgvar.SetAK(string(akBuf))
}
