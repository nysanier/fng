package pkgclient

import (
	"fmt"
	"log"
	"sync"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgvar"
)

const (
	OtsLimitOnce = 100 // 单次操作的数量, TODO: 测试环境=2，正式环境=100
)

// 封装ots的操作，对外提供方便使用的方法
type OtsClient struct {
	client *tablestore.TableStoreClient
}

var otsClientOnce sync.Once
var otsClient *OtsClient

// 单例模式
func GetOtsClient() *OtsClient {
	otsClientOnce.Do(func() {
		instanceName := "a3927top"
		endpoint := fmt.Sprintf("https://%v.%v.ots.aliyuncs.com", instanceName, pkgvar.RegionID)
		otsClient = &OtsClient{
			client: tablestore.NewClient(endpoint, instanceName, pkgvar.GetAK(), pkgvar.GetSK()),
		}
	})

	return otsClient
}

// 从ots中解析出来，或者格式到ots

// 返回的map已经初始化
func ParseColumn(cols []*tablestore.AttributeColumn) map[string]interface{} {
	valMap := map[string]interface{}{}
	for _, col := range cols {
		valMap[col.ColumnName] = col.Value
	}

	return valMap
}

func ParsePk(pk *tablestore.PrimaryKey) *OtsPks {
	pks := &OtsPks{}
	for _, col := range pk.PrimaryKeys {
		pks.PkList = append(pks.PkList, col.ColumnName)
		pks.ValList = append(pks.ValList, col.Value)
	}

	return pks
}

// 返回的map已经初始化
func ParseRow(row *tablestore.Row) map[string]interface{} {
	return ParseColumn(row.Columns)
}

// 返回的map已经初始化
func ParseRows(rows []*tablestore.Row) ([]*OtsPks, []map[string]interface{}) {
	var pksList []*OtsPks
	var valList []map[string]interface{}
	for _, row := range rows {
		pksList = append(pksList, ParsePk(row.PrimaryKey))
		valList = append(valList, ParseRow(row))
	}
	return pksList, valList
}

func FormatPk2(startPks, endPks *OtsPks) (*tablestore.PrimaryKey, *tablestore.PrimaryKey, error) {
	if err := startPks.CheckWithPks(endPks); err != nil {
		return nil, nil, err
	}

	startPk := new(tablestore.PrimaryKey)
	endPk := new(tablestore.PrimaryKey)
	for i, startPksPk := range startPks.PkList {
		startPksVal := startPks.ValList[i]
		endPksPk := endPks.PkList[i]
		endPksVal := endPks.ValList[i]

		// startPks中的空值(nil), 表示min, 否则用指定的val
		// endPks中的空值(nil), 表示max, 否则用指定的val
		if startPksVal == nil {
			startPk.AddPrimaryKeyColumnWithMinValue(startPksPk)
		} else {
			startPk.AddPrimaryKeyColumn(startPksPk, startPksVal)
		}

		if endPksVal == nil {
			endPk.AddPrimaryKeyColumnWithMaxValue(endPksPk)
		} else {
			endPk.AddPrimaryKeyColumn(endPksPk, endPksVal)
		}
	}

	return startPk, endPk, nil
}

// GetRow 类
// return columnNameValueMap, err
// err=nil时, map一定非空
func (p *OtsClient) GetRow(pkList []*tablestore.PrimaryKeyColumn, tableName string) (*tablestore.Row, error) {
	if len(pkList) == 0 {
		return nil, fmt.Errorf("invalid pkList(%v)", pkgfunc.FormatJson(pkList))
	}

	req := new(tablestore.GetRowRequest)
	criteria := new(tablestore.SingleRowQueryCriteria)
	criteria.PrimaryKey = &tablestore.PrimaryKey{PrimaryKeys: pkList}
	req.SingleRowQueryCriteria = criteria
	req.SingleRowQueryCriteria.TableName = tableName
	req.SingleRowQueryCriteria.MaxVersion = 1
	resp, err := otsClient.client.GetRow(req)
	if err != nil {
		log.Printf("otsClient.GetRow fail, err=%v", err)
		return nil, err
	}
	//log.Printf("resp: %v", pkgfunc.FormatJson(resp))

	row := &tablestore.Row{
		PrimaryKey: &resp.PrimaryKey,
		Columns:    resp.Columns,
	}

	return row, nil
}

type OtsPks struct {
	PkList  []string
	ValList []interface{}
}

func (p *OtsPks) Check() error {
	pkListLen := len(p.PkList)
	valListLen := len(p.ValList)

	if pkListLen != valListLen {
		return fmt.Errorf("pkListLen(%v) != valListLen(%v)", pkListLen, valListLen)
	}

	return nil
}

func (p *OtsPks) CheckWithPks(rhs *OtsPks) error {
	if err := p.Check(); err != nil {
		return fmt.Errorf("this Pks invalid, %v", err.Error())
	}

	if err := rhs.Check(); err != nil {
		return fmt.Errorf("that Pks invalid, %v", err.Error())
	}

	thisLen := len(p.PkList)
	thatLen := len(rhs.PkList)
	if thisLen != thatLen {
		return fmt.Errorf("thisLen(%v) != thatLen(%v)", thisLen, thatLen)
	}

	return nil
}

func (p *OtsClient) GetRowByPk(pk string, pkv interface{}, tableName string) (map[string]interface{}, error) {
	pkList := []*tablestore.PrimaryKeyColumn{
		{ColumnName: pk, Value: pkv},
	}

	row, err := p.GetRow(pkList, tableName)
	if err != nil {
		return nil, err
	}

	ret := ParseRow(row)
	return ret, nil
}

func (p *OtsClient) GetRowByPks(pks *OtsPks, tableName string) (map[string]interface{}, error) {
	if err := pks.Check(); err != nil {
		return nil, err
	}

	pkList := []*tablestore.PrimaryKeyColumn{}
	for i, pk := range pks.PkList {
		pkList = append(pkList, &tablestore.PrimaryKeyColumn{ColumnName: pk, Value: pks.ValList[i]})
	}

	row, err := p.GetRow(pkList, tableName)
	if err != nil {
		return nil, err
	}

	ret := ParseRow(row)
	return ret, nil
}

const (
	Order_ASC  = tablestore.FORWARD  // 升序
	Order_DESC = tablestore.BACKWARD // 降序
)

// GetRange 类, limit<1表示全部获取，每次处理100条, 会返回nextPrimary,
func (p *OtsClient) GetRange(startPk, endPk *tablestore.PrimaryKey, order tablestore.Direction, limit int, tableName string) ([]*tablestore.Row, *tablestore.PrimaryKey, error) {
	var rows []*tablestore.Row
	nextPk := startPk
	for {
		rowsOnce, nextPk2, err := p.getRangeInternal(nextPk, endPk, order, OtsLimitOnce, tableName)
		if err != nil {
			log.Printf("otsClient.getRangeInternal fail, err=%v", err)
			return nil, nil, err
		}

		rows = append(rows, rowsOnce...)
		nextPk = nextPk2

		// 已经达到limit的量了
		if limit > 0 && len(rows) > limit {
			break
		}

		// 没有下一个了
		if nextPk2 == nil {
			break
		}
	}

	return rows, nextPk, nil
}

func (p *OtsClient) getRangeInternal(startPk, endPk *tablestore.PrimaryKey, order tablestore.Direction, limit int, tableName string) ([]*tablestore.Row, *tablestore.PrimaryKey, error) {
	startPkLen := len(startPk.PrimaryKeys)
	endPkLen := len(endPk.PrimaryKeys)
	if startPkLen < 1 || startPkLen > 4 || startPkLen != endPkLen { // ots的pk最多4列
		return nil, nil, fmt.Errorf("invalid startPk=%v, endPk=%v", pkgfunc.FormatJson(startPk), pkgfunc.FormatJson(endPk))
	}

	req := new(tablestore.GetRangeRequest)
	criteria := new(tablestore.RangeRowQueryCriteria)
	criteria.StartPrimaryKey = startPk
	criteria.EndPrimaryKey = endPk
	criteria.Direction = order
	criteria.Limit = int32(limit)
	criteria.MaxVersion = 1 // 不使用多版本

	req.RangeRowQueryCriteria = criteria
	req.RangeRowQueryCriteria.TableName = tableName
	req.RangeRowQueryCriteria.MaxVersion = 1
	resp, err := otsClient.client.GetRange(req)
	if err != nil {
		log.Printf("otsClient.GetRange fail, err=%v", err)
		return nil, nil, err
	}
	//log.Printf("resp: %v", pkgfunc.FormatJson(resp))

	return resp.Rows, resp.NextStartPrimaryKey, nil
}

func (p *OtsClient) GetRangeAll(startPk, endPk *tablestore.PrimaryKey, tableName string) ([]*tablestore.Row, error) {
	order := Order_ASC // 顺序无关
	limit := 0
	rows, nextPk, err := p.GetRange(startPk, endPk, order, limit, tableName)
	if err != nil {
		log.Printf("GetRange fail, err=%v", err)
		return nil, err
	}

	_ = nextPk // 一定为nil
	return rows, err
}

func (p *OtsClient) GetRangeAllWithPks(startPks, endPks *OtsPks, tableName string) ([]*OtsPks, []map[string]interface{}, error) {
	startPk, endPk, err := FormatPk2(startPks, endPks)
	if err != nil {
		log.Printf("GetRangeAllWithPks fail, err=%v", err)
		return nil, nil, err
	}

	rows, err := p.GetRangeAll(startPk, endPk, tableName)
	if err != nil {
		log.Printf("GetRangeAll fail, err=%v", err)
		return nil, nil, err
	}

	pksList, valList := ParseRows(rows)
	return pksList, valList, nil
}

func (p *OtsClient) GetRangeWithPks(startPks, endPks *OtsPks, order tablestore.Direction, limit int, tableName string) ([]*tablestore.Row, *OtsPks, error) {
	startPk, endPk, err := FormatPk2(startPks, endPks)
	if err != nil {
		log.Printf("FormatPk2 fail, err=%v", err)
		return nil, nil, err
	}

	// 将nextPrimary转化为OtsPks格式
	rows, nextPrimary, err := p.GetRange(startPk, endPk, order, limit, tableName)
	if err != nil {
		log.Printf("GetRange fail, err=%v", err)
		return nil, nil, err
	}

	if nextPrimary == nil {
		return rows, nil, nil
	}

	nextPks := &OtsPks{}
	for _, pk := range nextPrimary.PrimaryKeys {
		nextPks.PkList = append(nextPks.PkList, pk.ColumnName)
		nextPks.ValList = append(nextPks.ValList, pk.Value)
	}
	return rows, nextPks, nil
}
