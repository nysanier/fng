package pkgconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgclient"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

// 定时加载配置(ots)
var (
	configLock      sync.RWMutex
	configMap       = make(map[string]map[string]interface{})
	configTableName = "fng_config" // 所有环境暂时都使用同一个表，通过block来区分env
	configTimer     *pkgfunc.Timer
)

// configMap的key格式，移除了dev/daily等前缀
func FormatConfigKey(block, section string) string {
	return fmt.Sprintf("%v#%v/%v", pkgvar.FnEnv, block, section)
}

// configMap的key格式，从config表加载过来的block和section已经包含了dev/daily前前缀，因此这里不额外添加
func FormatConfigKeyForSave(block, section string) string {
	return fmt.Sprintf("%v/%v", block, section)
}

// 每个env加载相关的所有配置项
func loadConfig() error {
	envStart := fmt.Sprintf("%v", pkgvar.FnEnv) // 比如`dev#x`一定是在`dev`和`dev~`之间的
	envEnd := fmt.Sprintf("%v~", pkgvar.FnEnv)
	startPks := &pkgclient.OtsPks{
		PkList:  []string{pkgvar.ConfigFieldPk1, pkgvar.ConfigFieldPk2},
		ValList: []interface{}{envStart, nil},
	}
	endPks := &pkgclient.OtsPks{
		PkList:  []string{pkgvar.ConfigFieldPk1, pkgvar.ConfigFieldPk2},
		ValList: []interface{}{envEnd, nil},
	}
	pksList, valList, err := pkgclient.GetOtsClient().GetRangeAllWithPks(startPks, endPks, configTableName)
	if err != nil {
		log.Printf("pkgclient.GetRangeAllWithPks fail, err=%v", err)
		return err
	}

	for i, pks := range pksList {
		itemMap := valList[i]
		//log.Printf("value=%v", itemMap[pkgvar.ConfigFieldValue])
		block := pks.ValList[0].(string)   // 0就是block
		section := pks.ValList[1].(string) // 1就是section
		if err := saveConfig(block, section, itemMap[pkgvar.ConfigFieldValue]); err != nil {
			log.Printf("saveConfig fail, err=%v", err)
			return err
		}
	}

	return nil
}

// TODO: 加锁又可以优化，但暂时先这样了
func saveConfig(block, section string, v interface{}) error {
	// 通过value这样一个json格式更通用，因为mysql等rds作为配制源的话，扩展配置字段没有那么方便，且ots行不会拉的那么长！
	value := "null"
	if v != nil {
		value = v.(string)
	}

	itemMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &itemMap); err != nil {
		log.Printf("json.Unmarshal fail, err=%v, ignore the value=%v", err, value)
		return nil
	}

	key := FormatConfigKeyForSave(block, section)

	// 最小化加锁
	configLock.Lock()
	configMap[key] = itemMap
	configLock.Unlock()

	return nil
}

//  每个env只加载相关的dev#base/common这一个配置项
func loadOneConfig() error {
	block := "base"
	section := "common"
	pk1 := fmt.Sprintf("%v#%v", pkgvar.FnEnv, block)
	pk2 := section
	pks := &pkgclient.OtsPks{
		PkList:  []string{pkgvar.ConfigFieldPk1, pkgvar.ConfigFieldPk2},
		ValList: []interface{}{pk1, pk2},
	}

	itemMap, err := pkgclient.GetOtsClient().GetRowByPks(pks, configTableName)
	if err != nil {
		log.Printf("pkgclient.GetRowByPks fail, err=%v", err)
		return err
	}

	// 通过value这样一个json格式更通用，因为mysql等rds作为配制源的话，扩展配置字段没有那么方便，且ots行不会拉的那么长！
	if err := saveConfig(block, section, itemMap[pkgvar.ConfigFieldValue]); err != nil {
		log.Printf("saveConfig fail, err=%v", err)
		return err
	}

	log.Printf("loadOneConfig ok")
	return nil
}

// 每5分钟更新一次配置
func StartConfigUpdater() {
	interval := 300
	configTimer = pkgfunc.NewTimer(loadConfig, time.Second*time.Duration(interval))
	configTimer.Start()
	log.Printf("start config updater ok, interval=%v", interval)
}

// 只支持获取基础类型的配置，即string,int,bool等
func GetConfig(block, section, item string) interface{} {
	key := FormatConfigKey(block, section)

	// 最小化加锁
	configLock.RLock()
	itemMap := configMap[key]
	configLock.RUnlock()

	return itemMap[item]
}

func GetConfigString(block, section, item string) string {
	v := GetConfig(block, section, item)
	v2 := v.(string)
	return v2
}

func GetConfigStringWithDefault(block, section, item string, defVal string) string {
	v := GetConfig(block, section, item)
	if v2, ok := v.(string); ok {
		return v2
	}

	return defVal
}

// 32位整型
func GetConfigInteger(block, section, item string) int {
	v := GetConfig(block, section, item)
	v2 := v.(float64)
	return int(v2)
}

func GetConfigIntegerWithDefault(block, section, item string, defVal int) int {
	v := GetConfig(block, section, item)
	if v2, ok := v.(float64); ok {
		return int(v2)
	}

	return defVal
}

func GetConfigBoolean(block, section, item string) bool {
	v := GetConfig(block, section, item)
	v2 := v.(bool)
	return v2
}

func GetConfigBooleanWithDefault(block, section, item string, defVal bool) bool {
	v := GetConfig(block, section, item)
	if v2, ok := v.(bool); ok {
		return v2
	}

	return defVal
}
