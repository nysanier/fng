package fnmetrics

import (
	"testing"
	"time"
)

func Test_Metrics(t *testing.T) {
	metrics := NewMetrics("TEST")
	metrics.StartMetrics()
	defer metrics.Close()

	// 最多运行20秒
	for i := 0; i < 33*20; i++ {
		time.Sleep(time.Millisecond * 30)
		if i < 100 {
			metrics.AddValue(32 * 1024)
		}
	}

	time.Sleep(time.Second)
}

/*
=== RUN   Test_Metrics
2022/04/25 00:42:05 <TEST> avgSpeed: 1.07MiB/s, curSpeed: 1.04MiB/s
2022/04/25 00:42:06 <TEST> avgSpeed: 1.04MiB/s, curSpeed: 1.03MiB/s
2022/04/25 00:42:07 <TEST> avgSpeed: 1.04MiB/s, curSpeed: 1.03MiB/s
2022/04/25 00:42:08 <TEST> avgSpeed: 805.46KiB/s, curSpeed: 799.13KiB/s
2022/04/25 00:42:09 <TEST> avgSpeed: 646.49KiB/s, curSpeed: 642.41KiB/s
2022/04/25 00:42:10 <TEST> avgSpeed: 531.09KiB/s, curSpeed: 528.33KiB/s
2022/04/25 00:42:11 <TEST> avgSpeed: 450.84KiB/s, curSpeed: 448.85KiB/s
2022/04/25 00:42:12 <TEST> avgSpeed: 396.19KiB/s, curSpeed: 394.65KiB/s
2022/04/25 00:42:13 <TEST> avgSpeed: 349.77KiB/s, curSpeed: 348.57KiB/s
2022/04/25 00:42:15 <TEST> avgSpeed: 313.03KiB/s, curSpeed: 312.07KiB/s
2022/04/25 00:42:16 <TEST> avgSpeed: 285.73KiB/s, curSpeed: 284.93KiB/s
2022/04/25 00:42:17 <TEST> avgSpeed: 262.81KiB/s, curSpeed: 262.13KiB/s
2022/04/25 00:42:18 <TEST> avgSpeed: 241.51KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:19 <TEST> avgSpeed: 224.86KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:20 <TEST> avgSpeed: 209.09KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:21 <TEST> avgSpeed: 195.40KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:22 <TEST> avgSpeed: 183.38KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:23 <TEST> avgSpeed: 173.66KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:24 <TEST> avgSpeed: 164.91KiB/s, curSpeed: 0.00B/s
2022/04/25 00:42:25 <TEST> avgSpeed: 156.27KiB/s, curSpeed: 0.00B/s
--- PASS: Test_Metrics (21.31s)
2022/04/25 00:42:25 <TEST> totalSize: 3.28MiB, totalCount: 281, avgSpeed: 154.01KiB/s
PASS
coverage: 89.9% of statements in ../../../fngo/...
*/
