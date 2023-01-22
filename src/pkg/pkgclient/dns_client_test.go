package pkgclient

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetDnsClient(t *testing.T) {
	patches := NewPatches()

	Convey("Test_GetDnsClient", t, func() {
		defer patches.Reset()

		Convey("ok", func() {
			defer patches.Reset()
			
			//ListA3927Dns()
			SetA3927Dns("test", "0.0.0.3")
		})
	})
}
