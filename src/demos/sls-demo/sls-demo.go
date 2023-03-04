package main

import (
	"fmt"
	"log"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/nysanier/fng/src/pkg/pkgclient"
	"github.com/nysanier/fng/src/pkg/pkgenv"
	"github.com/nysanier/fng/src/pkg/pkglog"
	"github.com/nysanier/fng/src/pkg/pkglog/logimpl"
)

var (

	// endpoint参考 https://help.aliyun.com/document_detail/29008.html
	// 内网地址 cn-hangzhou-intranet.log.aliyuncs.com
	Endpoint     = "cn-hangzhou.log.aliyuncs.com" // 公网地址
	AK           = os.Getenv("AK")
	SK           = os.Getenv("SK")
	ProjectName  = "a3927top"
	LogStoreName = "a3927top-nginx"
)

func main() {
	//producerConfig := producer.GetDefaultProducerConfig()
	//producerConfig.Endpoint = "cn-hangzhou.log.aliyuncs.com" // 公网地址
	//producerConfig.AccessKeyID = os.Getenv("AK")
	//producerConfig.AccessKeySecret = os.Getenv("SK")
	//producerInstance := producer.InitProducer(producerConfig)
	//ch := make(chan os.Signal)
	////signal.Notify(ch, os.Kill, os.Interrupt)
	//signal.Notify(ch)
	//producerInstance.Start() // 启动producer实例
	//
	//sig := <-ch
	//log.Printf("sig: %v", sig)
	writeLog2()

	log.Printf("ok")
}

func writeLog2() {
	pkgenv.LoadEnv()
	logimpl.InitSlsLog()
	pkglog.Infov("Evt1", "k1", "v1", "k2", "v2")
	pkglog.Warnv("EvtXXXFail", "Error", "ErrYYY")
	pkglog.Errorv("Evt3")
}

func writeLog() {
	// 创建日志服务Client。
	client := pkgclient.GetSlsClient()

	// 向Logstore写入数据。
	logs := []*sls.Log{}
	for logIdx := 0; logIdx < 2; logIdx++ { // 2条日志
		content := []*sls.LogContent{}
		for colIdx := 0; colIdx < 3; colIdx++ { // 每条日志3个字段
			content = append(content, &sls.LogContent{
				Key:   proto.String(fmt.Sprintf("col_%d", colIdx)), // key名字
				Value: proto.String(fmt.Sprintf("%d", colIdx*10)),
			})
		}
		log := &sls.Log{
			Time:     proto.Uint32(uint32(time.Now().Unix())),
			Contents: content,
		}
		logs = append(logs, log)
	}

	logGroup := &sls.LogGroup{
		Topic:  proto.String("writeLogWithGolang"),
		Source: proto.String("1.2.3.4"),
		Logs:   logs,
	}

	if err := client.PutLogs(ProjectName, LogStoreName, logGroup); err != nil {
		log.Fatalf("PutLogs failed %v", err)
		os.Exit(1)
	}
}

/*
__source__:1.2.3.4
__tag__:__receive_time__:1677935184
__topic__:writeLogWithGolang
__tag__:__client_ip__:183.159.122.125
col_0:0
col_1:10
col_2:20

... 重复3遍
*/
