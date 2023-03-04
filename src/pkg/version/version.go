package version

import (
	"strconv"
	"time"

	"github.com/nysanier/fng/src/pkg/pkgfunc"
	"github.com/nysanier/fng/src/pkg/pkgutil"
)

// 编译时通过ldflags注入
var (
	AppVer    = "xx"
	GitCommit = "yy"
	BuildTime = "zz"
)

const (
	ShortGitCommitLength = 8
)

func GetShortGitCommit() string {
	if len(GitCommit) < ShortGitCommitLength {
		return GitCommit
	}

	return GitCommit[:ShortGitCommitLength]
}

func GetBuildTimeStr() string {
	v, _ := strconv.ParseInt(BuildTime, 10, 64)
	t := time.Unix(v, 0)
	t2 := pkgutil.ToCstTime(t)
	str := pkgfunc.GetRFC3339TimeStr(t2)
	return str
}

func GetAppVersion() string {
	return AppVer
}
