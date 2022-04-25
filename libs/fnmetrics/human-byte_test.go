package fnmetrics

import (
	"log"
	"testing"
)

func Test_HumanByte(t *testing.T) {
	b := HumanByte(123)
	log.Printf("%v", b)
	kb := HumanByte(12_123)
	log.Printf("%v", kb)
	mb := HumanByte(12_123456)
	log.Printf("%v", mb)
	gb := HumanByte(12_123456789)
	log.Printf("%v", gb)
	tb := HumanByte(12345_123456789)
	log.Printf("%v", tb)
	pb := HumanByte(12_123456_123456789)
	log.Printf("%v", pb)
	pb2 := HumanByte(12345_123456_123456789)
	log.Printf("%v", pb2)
}

/*
=== RUN   Test_HumanByte
2022/04/25 00:45:05 123.00B
2022/04/25 00:45:05 12.12KiB
2022/04/25 00:45:05 12.12MiB
2022/04/25 00:45:05 12.12GiB
2022/04/25 00:45:05 12.35TiB
2022/04/25 00:45:05 12.12PiB
2022/04/25 00:45:05 12345.12PiB
--- PASS: Test_HumanByte (0.00s)
PASS
coverage: 11.6% of statements in ../../../fngo/...

*/
