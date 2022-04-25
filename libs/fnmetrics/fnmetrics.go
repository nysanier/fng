package fnmetrics

import (
	"log"
	"time"
)

const (
	Cnt         = 3    // 阶段个数
	ChanCap     = 1024 // channel 容量, 用于暂存packet
	PacketLimit = 100  // 每个包分组中包个数(packet count)上限
)

// 当前速度算法: 当前总量/当前总时间, PacketLimit ~ 2*PacketLimit 个包的平均速度作为
//   当前总量: 最后2个阶段的总量
//   当前总时间: 最后2个阶段的总时间
// 平均速度算法: 总量/总时间

type metricsStruct struct {
	CurSize int64
	CurTime int64
}

type Metrics struct {
	startTime        int64      // 第一个包的开始时间
	totalSize        int64      // 包总量
	totalCount       int64      // 包总数
	lastStartTime    [Cnt]int64 // 阶段x最后一个包的开始时间
	lastTotalSize    [Cnt]int64 // 阶段x包总量
	lastTotalCount   [Cnt]int64 // 阶段x包总数, [Cnt-1]-[Cnt-2]<=PacketLimit, 也就是前后elem差值不会超过PacketLimit个包
	ch               chan metricsStruct
	lastDumpTime     int64
	prefix           string
	lastAddValueTime int64
	ticker           *time.Ticker
}

func NewMetrics(prefix string) *Metrics {
	p := &Metrics{
		ch:     make(chan metricsStruct, ChanCap),
		prefix: prefix,
	}

	curTime := time.Now().UnixNano()
	for i := 0; i < Cnt; i++ {
		p.lastStartTime[i] = curTime
	}

	return p
}

// 如果连续n秒没有数据, 那么往chan中放入一个0值
func (p *Metrics) AddValue(n int) {
	p.ch <- metricsStruct{int64(n), time.Now().UnixNano()}
	p.lastAddValueTime = time.Now().UnixNano()
}

func (p *Metrics) Close() {
	close(p.ch)
	if p.ticker != nil {
		p.ticker.Stop()
	}
}

func (p *Metrics) StartMetrics() {
	// 每100ms产生一个空数据块, 1s有10个空数据块, 那么PacketLimit/10秒后速度归零
	p.ticker = time.NewTicker(time.Millisecond * 100)

	go func() {
		for {
			select {
			case _, ok := <-p.ticker.C:
				// 定时器被关闭了, 退出逻辑
				if !ok {
					return
				}

				// 如果有100ms没有更新数据了, 那么在chan中翻入一个0值
				if time.Now().UnixNano()-p.lastAddValueTime > 1e6*100 {
					p.ch <- metricsStruct{0, time.Now().UnixNano()}
				}
			}
		}
	}()

	go p.doMetrics()
}

func (p *Metrics) doMetrics() {
	for {
		select {
		case st, ok := <-p.ch:
			if !ok {
				curTime := time.Now().UnixNano()
				avgTotalSize := p.totalSize         // Byte
				avgTotalNs := curTime - p.startTime // ns
				if avgTotalNs == 0 {
					avgTotalNs = 1
				}
				avgSpeed := float64(avgTotalSize) * 1e9 / float64(avgTotalNs) // byte/ns = byte*1e9/s
				log.Printf("<%v> totalSize: %v, totalCount: %v, avgSpeed: %v/s", p.prefix, HumanByte(float64(p.totalSize)), p.totalCount, HumanByte(avgSpeed))
				return
			}

			curTime := st.CurTime
			curSize := st.CurSize

			if p.startTime == 0 {
				p.startTime = curTime
				p.lastDumpTime = curTime
			}

			p.totalCount += 1
			p.totalSize += curSize
			// 如果最新阶段的包个数达到上限了, 那么丢弃最老的阶段数据, 创建一个新的阶段, 并将所有阶段往前移动一个
			if p.lastTotalCount[Cnt-1] >= p.lastTotalCount[Cnt-2]+PacketLimit {
				for i := 0; i < Cnt-1; i++ {
					p.lastStartTime[i] = p.lastStartTime[i+1]
					p.lastTotalSize[i] = p.lastTotalSize[i+1]
					p.lastTotalCount[i] = p.lastTotalCount[i+1]
				}
			}
			p.lastStartTime[Cnt-1] = curTime
			p.lastTotalSize[Cnt-1] = p.totalSize
			p.lastTotalCount[Cnt-1] = p.totalCount

			// 每秒最多打印一次
			if curTime-p.lastDumpTime > 1e9 {
				curTotalSize := p.lastTotalSize[Cnt-1] - p.lastTotalSize[0] // Byte
				curTotalNs := p.lastStartTime[Cnt-1] - p.lastStartTime[0]   // ns
				if curTotalNs == 0 {
					curTotalNs = 1
				}
				curSpeed := float64(curTotalSize) * 1e9 / float64(curTotalNs)

				avgTotalSize := p.totalSize
				avgTotalNs := curTime - p.startTime
				if avgTotalNs == 0 {
					avgTotalNs = 1
				}
				avgSpeed := float64(avgTotalSize) * 1e9 / float64(avgTotalNs)
				log.Printf("<%v> avgSpeed: %v/s, curSpeed: %v/s", p.prefix, HumanByte(avgSpeed), HumanByte(curSpeed))
				p.lastDumpTime = curTime
			}
		}
	}
}
