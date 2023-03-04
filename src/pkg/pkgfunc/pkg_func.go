package pkgfunc

import (
	"math/rand"
	"time"
)

func InitRand() {
	rand.Seed(time.Now().UnixNano())
}
