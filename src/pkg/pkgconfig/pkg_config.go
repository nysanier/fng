package pkgconfig

import (
	"fmt"
	"log"

	"github.com/nysanier/fng/src/pkg/pkgclient"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

// TODO: 启动一个协程定时加载配置
// 加载配置(ots)
func LoadConfig() {
	blockVal := fmt.Sprintf("fng_%v", pkgvar.FnEnv)
	otsPk := pkgclient.OtsPk{
		PkList:  []string{"block", "section"},
		ValList: []interface{}{blockVal, "common"},
	}
	tableName := "config"
	otsRow, err := pkgclient.GetOtsRow2(otsPk, tableName)
	if err != nil {
		log.Printf("pkgclient.GetOtsRow fail, err=%v", err)
		return
	}

	log.Printf("otsRow=%v", pkgfunc.FormatJson(otsRow))
	if otsRow[pkgvar.ConfigValue] == nil {
		log.Printf("config value is nil, otsPk=%v", pkgfunc.FormatJson(otsPk))
		panic("config value is nil")
		return
	}
}

//func GetConfig(block, section, item string) interface{} {
//
//}

//func GetConfigString(block, section, item string) string {
//	v := GetConfig(block, section, item)
//	ret := v.(string)
//	return ret
//}
