package main

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"time"

	"github.com/nysanier/fng/libs/fnmetrics"
	uuid "github.com/satori/go.uuid"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	sha1Str := getSha1Str(12345) // macos上需要运行3s
	log.Printf(sha1Str)

	time.Sleep(time.Second)
}

func getUuid() []byte {
	uuidArr := uuid.NewV4()
	return uuidArr[:]
}

func getSha1Str(n int) string {
	metrics := fnmetrics.NewMetrics("SHA1")
	metrics.StartMetrics()
	defer metrics.Close()

	h := sha1.New()

	// 先加入一个uuid, 保证每次计算的data都是不同的, 排除缓存等干扰因素
	buf := getUuid()
	metrics.AddValue(len(buf))
	if _, err := h.Write(buf); err != nil {
		panic(err)
	}

	// 单包>=32K的时候计算效率最高, 这个值也是网络包的大小
	buf2 := make([]byte, 256*1024) // <SHA1> avgSpeed: 1.04GiB/s, curSpeed: 1.07GiB/s
	// buf2 := make([]byte, 32*1024) // <SHA1> avgSpeed: 976.12MiB/s, curSpeed: 1.04GiB/s
	// buf2 := make([]byte, 4096) // <SHA1> avgSpeed: 871.42MiB/s, curSpeed: 941.83MiB/s
	// buf2 := make([]byte, 1024) // <SHA1> avgSpeed: 637.08MiB/s, curSpeed: 648.29MiB/s
	// buf2 := make([]byte, uuid.Size) // <SHA1> avgSpeed: 42.08MiB/s, curSpeed: 45.54MiB/s
	// buf2 := getUuid() // <SHA1> avgSpeed: 42.43MiB/s, curSpeed: 45.82MiB/s

	for i := 0; i < n; i++ {
		// buf2 := make([]byte, uuid.Size) // <SHA1> avgSpeed: 37.78MiB/s, curSpeed: 38.40MiB/s
		// buf2 := getUuid() // <SHA1> avgSpeed: 22.79MiB/s, curSpeed: 24.83MiB/s
		metrics.AddValue(len(buf2))
		if _, err := h.Write(buf2); err != nil {
			panic(err)
		}
	}

	r := h.Sum(nil)
	str := hex.EncodeToString(r)
	return str
}
