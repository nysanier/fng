package pkgconf

import (
	"fmt"
	"sync"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkglog"
)

var (
	configLock sync.RWMutex
	// 比如 base.common = {"x"=1,"y"="a"}, 即 base.common.x = 1
	configMap   = make(map[string]map[string]interface{})
	configTimer *pkgfunc.Timer
)

// 每5分钟更新一次配置
func InitConf(loadConfig pkgfunc.TimerFunc) {
	interval := 300
	configTimer = pkgfunc.NewTimer(loadConfig, time.Second*time.Duration(interval))
	configTimer.Start()
	pkglog.Infov("EvtConfInitOK",
		"UpdateInterval", interval)
}

// configMap的key格式，比如 base.common
func FormatConfigKey(block, section string) string {
	return fmt.Sprintf("%v.%v", block, section)
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

func SetItemMap(block, section string, itemMap map[string]interface{}) {
	key := FormatConfigKey(block, section)

	// 最小化加锁
	configLock.Lock()
	configMap[key] = itemMap
	configLock.Unlock()
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
