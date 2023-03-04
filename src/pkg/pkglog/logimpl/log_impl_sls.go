package logimpl

import (
	"log"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/nysanier/fng/src/pkg/pkgclient"
)

const (
	ProjectName  = "a3927top"
	LogStoreName = "a3927top-nginx"
)

func NewLogImplSls() *LogImplSls {
	return &LogImplSls{}
}

type LogImplSls struct {
}

func (p *LogImplSls) WriteLog(keys []string, vals []string) {
	logs := []*sls.Log{}
	content := []*sls.LogContent{}
	for i, key := range keys {
		content = append(content, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(vals[i]),
		})
	}
	slsLog := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content, // 多对kv
	}
	logs = append(logs, slsLog)

	logGroup := &sls.LogGroup{
		//Topic:  proto.String("writeLogWithGolang"),
		//Source: proto.String("1.2.3.4"),
		Logs: logs, // 1条日志
	}

	if err := pkgclient.GetSlsClient().PutLogs(ProjectName, LogStoreName, logGroup); err != nil {
		log.Printf("EvtPutLogsFail, Error=%v", err) // 初始化还未成功，因此不能循环调用，先打印到本地
	}
}
