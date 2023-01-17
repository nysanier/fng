package version

// 编译时通过ldflags注入
var (
	AppVer    = "xx"
	GitCommit = "yy"
	BuildTime = "zz"
)

func GetShortGitCommit() string {
	if len(GitCommit) < 10 {
		return GitCommit
	}

	return GitCommit[:10]
}
