package version

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
