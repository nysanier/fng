package fnmetrics

import "fmt"

const (
	Bi  = int64(1)
	KiB = Bi * 1000
	MiB = KiB * 1000
	GiB = MiB * 1000
	TiB = GiB * 1000
	PiB = TiB * 1000

	KB = Bi * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
)

var (
	byteList = []int64{Bi, KiB, MiB, GiB, TiB, PiB}
	nameList = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
)

// 除了PiB, 其他情况下最多返回3位数字
func HumanByteInt64(n int64) string {
	return HumanByte(float64(n))
}

func HumanByte(n float64) string {
	byteListLen := len(byteList)
	nameListLen := len(nameList)
	if byteListLen != nameListLen {
		panic(fmt.Sprintf("byteListLen(%v) != nameListLen(%v)", byteListLen, nameListLen))
	}

	for i, v := range byteList {
		if n < float64(v) {
			// Bi
			if i == 0 {
				return fmt.Sprintf("%.2f%v", n, nameList[i])
			}

			// KiB/MiB等
			return fmt.Sprintf("%.2f%v", n/float64(byteList[i-1]), nameList[i-1])
		}
	}

	// 当前PiB是兜底
	return fmt.Sprintf("%.2f%v", n/float64(PiB), "PiB")
}
