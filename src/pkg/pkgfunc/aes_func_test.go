package pkgfunc

import (
	"log"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/nysanier/fng/src/pkg/pkgvar"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_EncryptByAES(t *testing.T) {
	patches := NewPatches()
	//var data []byte

	Convey("Test_GetDnsClient", t, func() {
		defer patches.Reset()

		Convey("ok", func() {
			defer patches.Reset()

			pkgvar.AesKey = []byte("aaaaaaaaaaxxxxxx") // 补齐到16位
			//var err error
			str := "bbbbbbbbbbxxxxxx"
			r, err1 := EncryptByAES([]byte(str))
			buf, err2 := DecryptByAES(r)
			log.Printf("err1=%v, err2=%v", err1, err2)
			str2 := string(buf)
			So(str2, ShouldEqual, str)
		})
	})
}
