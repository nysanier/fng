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

// TODO: 启动一个协程定时加载配置
// 加载配置(ots)
func LoadConfig() {
	pk1 := fmt.Sprintf("fng_%v", pkgvar.FnEnv)
	pk2 := "common"
	otsPk := pkgclient.OtsPk{
		PkList:  []string{pkgvar.ConfigFieldPk1, pkgvar.ConfigFieldPk2},
		ValList: []interface{}{pk1, pk2},
	}
	otsRow, err := pkgclient.GetOtsRow2(otsPk, configTableName)
	if err != nil {
		log.Printf("pkgclient.GetOtsRow fail, err=%v", err)
		return
	}

	log.Printf("otsRow=%v", pkgfunc.FormatJson(otsRow))
	if otsRow[pkgvar.ConfigFieldValue] == nil {
		log.Printf("config value is nil, otsPk=%v", pkgfunc.FormatJson(otsPk))
		panic("config value is nil")
		return
	}
}

var (
	//cfgMap      sync.Map
	//config      Config
	configLock      sync.RWMutex
	configMap       = make(map[string]map[string]interface{})
	configTableName = "config"
	configTimer     *pkgfunc.Timer
)

// configMap的key格式，移除了fng_daily等前缀
func FormatConfigKey(block, section string) string {
	return fmt.Sprintf("%v#%v", block, section)
}

// TODO: 支持加载所有的配置项，目前仅支持加载fng_env#common配置项
func loadConfig() error {
	env := fmt.Sprintf("fng_%v", pkgvar.FnEnv)
	block := "base"
	section := "common"
	pk1 := fmt.Sprintf("%v#%v", env, block)
	pk2 := section
	otsPk := pkgclient.OtsPk{
		PkList:  []string{pkgvar.ConfigFieldPk1, pkgvar.ConfigFieldPk2},
		ValList: []interface{}{pk1, pk2},
	}
	otsRow, err := pkgclient.GetOtsRow2(otsPk, configTableName)
	if err != nil {
		log.Printf("pkgclient.GetOtsRow fail, err=%v", err)
		return err
	}

	value := "null"
	valIntf := otsRow[pkgvar.ConfigFieldValue]
	if valIntf != nil {
		value = valIntf.(string)
	}

	//log.Printf("value=%v", value)

	itemMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &itemMap); err != nil {
		log.Printf("json.Unmarshal fail, err=%v", err)
		return err
	}

	key := FormatConfigKey(block, section)

	// 最小化加锁
	configLock.Lock()
	configMap[key] = itemMap
	configLock.Unlock()

	log.Printf("loadConfig ok")
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
