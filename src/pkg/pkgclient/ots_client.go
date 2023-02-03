package pkgclient

import (
	"fmt"
	"log"
	"sync"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

var (
	otsClientOnce sync.Once
	otsClient     *tablestore.TableStoreClient
)

func GetOtsClient() *tablestore.TableStoreClient {
	instanceName := "a3927top"
	endpoint := fmt.Sprintf("https://%v.%v.ots.aliyuncs.com", instanceName, pkgvar.RegionID)

	otsClientOnce.Do(func() {
		otsClient = tablestore.NewClient(endpoint, instanceName, pkgvar.GetAK(), pkgvar.GetSK())
	})

	return otsClient
}

// return columnNameValueMap, err
// err=nil时, map一定非空
func GetOtsRow(pkList []*tablestore.PrimaryKeyColumn, tableName string) (map[string]interface{}, error) {
	if len(pkList) == 0 {
		return nil, fmt.Errorf("invalid pkList(%v)", pkgfunc.FormatJson(pkList))
	}

	client := GetOtsClient()

	req := new(tablestore.GetRowRequest)
	criteria := new(tablestore.SingleRowQueryCriteria)
	putPk := new(tablestore.PrimaryKey)
	putPk.PrimaryKeys = pkList
	criteria.PrimaryKey = putPk
	req.SingleRowQueryCriteria = criteria
	req.SingleRowQueryCriteria.TableName = tableName
	req.SingleRowQueryCriteria.MaxVersion = 1
	resp, err := client.GetRow(req)
	if err != nil {
		log.Printf("client.GetOtsRow fail, err=%v", err)
		return nil, err
	}
	//log.Printf("resp: %v", pkgfunc.FormatJson(resp))

	ret := map[string]interface{}{}
	for _, col := range resp.Columns {
		ret[col.ColumnName] = col.Value
	}

	return ret, nil
}

type OtsPk struct {
	PkList  []string
	ValList []interface{}
}

func GetOtsRow1(pk string, pkv interface{}, tableName string) (map[string]interface{}, error) {
	pkList := []*tablestore.PrimaryKeyColumn{
		{ColumnName: pk, Value: pkv},
	}

	return GetOtsRow(pkList, tableName)
}

func GetOtsRow2(otsPk OtsPk, tableName string) (map[string]interface{}, error) {
	l1 := len(otsPk.PkList)
	l2 := len(otsPk.ValList)
	if l1 > l2 { // 取两者较小者
		l1 = l2
	}

	pkList := []*tablestore.PrimaryKeyColumn{}
	for i, pk := range otsPk.PkList {
		pkList = append(pkList, &tablestore.PrimaryKeyColumn{ColumnName: pk, Value: otsPk.ValList[i]})
	}

	return GetOtsRow(pkList, tableName)
}
