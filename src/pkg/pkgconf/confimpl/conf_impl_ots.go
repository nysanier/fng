package confimpl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgclient"
	"github.com/nysanier/fng/src/pkg/pkgconf"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkglog"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

// 定时加载配置(ots)
const (
	tableName = "fng_config" // 所有环境暂时都使用同一个表，通过block来区分env

	// 配置表2个pk列和1个属性列
	ConfigFieldPkBlock   = "block"
	ConfigFieldPkSection = "section"
	ConfigFieldValue     = "value"
)

var (
	configTimer *pkgfunc.Timer
)

// 每5分钟更新一次配置
func StartConfigUpdater() {
	interval := 300
	configTimer = pkgfunc.NewTimer(loadConfig, time.Second*time.Duration(interval))
	configTimer.Start()
	pkglog.Infov("EvtConfStartUpdaterOK",
		"Interval", interval)
}

// 每个env加载相关的所有配置项
func loadConfig() error {
	envStart := fmt.Sprintf("%v", pkgvar.FnEnv) // 比如`dev#x`一定是在`dev`和`dev~`之间的
	envEnd := fmt.Sprintf("%v~", pkgvar.FnEnv)
	startPks := &pkgclient.OtsPks{
		PkList:  []string{ConfigFieldPkBlock, ConfigFieldPkSection},
		ValList: []interface{}{envStart, nil},
	}
	endPks := &pkgclient.OtsPks{
		PkList:  []string{ConfigFieldPkBlock, ConfigFieldPkSection},
		ValList: []interface{}{envEnd, nil},
	}
	pksList, valList, err := pkgclient.GetOtsClient().GetRangeAllWithPks(startPks, endPks, tableName)
	if err != nil {
		pkglog.Warnv("EvtOtsGetRangeAllWithPksFail",
			"Error", err)
		return err
	}

	var sectionCount int
	for i, pks := range pksList {
		itemMap := valList[i]
		//log.Printf("value=%v", itemMap[pkgvar.ConfigFieldValue])
		pkBlock := pks.ValList[0].(string) // 0就是block
		section := pks.ValList[1].(string) // 1就是section
		if err := saveConfig(pkBlock, section, itemMap[ConfigFieldValue]); err != nil {
			pkglog.Warnv("EvtOtsConfSaveConfigFailButContinue",
				"Error", err)
			continue
		}

		sectionCount += 1
	}

	pkglog.Infov("EvtOtsConfLoadConfigOK",
		"SectionCount", sectionCount)
	return nil
}

func saveConfig(pkBlock, section string, v interface{}) error {
	// 移除env前缀，比如pkBlock=daily#base，则block=base
	l := strings.Split(pkBlock, "#")
	if len(l) != 2 {
		return fmt.Errorf("invalid pkBlock(%v)", pkBlock)
	}
	block := l[1]

	// 通过value这样一个json格式更通用，因为mysql等rds作为配制源的话，扩展配置字段没有那么方便，且ots行不会拉的那么长！
	value := "null"
	if v != nil {
		value = v.(string)
	}

	itemMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &itemMap); err != nil {
		pkglog.Warnv("EvtJsonUnmarshalFail",
			"Error", err,
			"Value", value)
		return err
	}

	pkgconf.SetItemMap(block, section, itemMap)
	return nil
}

// configMap的key格式，移除了dev/daily等前缀
func FormatConfigKey(block, section string) string {
	return fmt.Sprintf("%v#%v/%v", pkgvar.FnEnv, block, section)
}

// configMap的key格式，从config表加载过来的block和section已经包含了dev/daily前前缀，因此这里不额外添加
func FormatConfigKeyForSave(block, section string) string {
	return fmt.Sprintf("%v/%v", block, section)
}

//  每个env只加载相关的dev#base/common这一个配置项
func loadOneConfig() error {
	block := "base"
	section := "common"
	pk1 := fmt.Sprintf("%v#%v", pkgvar.FnEnv, block)
	pk2 := section
	pks := &pkgclient.OtsPks{
		PkList:  []string{ConfigFieldPkBlock, ConfigFieldPkSection},
		ValList: []interface{}{pk1, pk2},
	}

	itemMap, err := pkgclient.GetOtsClient().GetRowByPks(pks, tableName)
	if err != nil {
		pkglog.Warnv("EvtOtsConfGetRowByPksFail",
			"Error", err)
		return err
	}

	// 通过value这样一个json格式更通用，因为mysql等rds作为配制源的话，扩展配置字段没有那么方便，且ots行不会拉的那么长！
	if err := saveConfig(block, section, itemMap[ConfigFieldValue]); err != nil {
		pkglog.Warnv("EvtOtsConfSaveConfigFail",
			"Error", err)
		return err
	}

	pkglog.Infov("EvtOtsConfLoadOneConfigOK")
	return nil
}
