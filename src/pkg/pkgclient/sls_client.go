package pkgclient

import (
	"fmt"
	"sync"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

// 封装ots的操作，对外提供方便使用的方法
type SlsClient struct {
	client *tablestore.TableStoreClient
}

var slsClientOnce sync.Once
var slsClient sls.ClientInterface

// 单例模式
func GetSlsClient() sls.ClientInterface {
	slsClientOnce.Do(func() {
		// 比如 cn-hangzhou.log.aliyuncs.com
		var Endpoint = fmt.Sprintf("%v.log.aliyuncs.com", pkgvar.GetRegionID()) // 公网地址
		//if local { // TODO: 暂时不考虑公网上下行流量
		//	Endpoint = fmt.Sprintf("%v-intranet.log.aliyuncs.com", pkgvar.RegionID) // 内网地址
		//}
		slsClient = sls.CreateNormalInterface(Endpoint, pkgvar.GetAK(), pkgvar.GetSK(), "")
	})

	return slsClient
}
