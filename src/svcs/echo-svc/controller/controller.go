package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nysanier/fng/src/pkg/pkgddns"
	"github.com/nysanier/fng/src/pkg/pkglog"
	"github.com/nysanier/fng/src/pkg/pkgutil"
	"github.com/nysanier/fng/src/pkg/version"
)

const (
	Fn3339     = "Mon, 02 Jan 2006 3:04:05 PM"
	BodyFormat = `hello, echo-svc,
    current is:  %v
    remote addr: %v

> app version: %v.%v
> build time:  %v
> start time:  %v
> service ip:  %v`
)

func Index(ctx *gin.Context) {
	curTimeStr := getFn1123CstTimeStr()
	remoteAddr := ctx.Request.RemoteAddr
	pkglog.Infov("EvtDumpRemoteAddr",
		"Addr", remoteAddr)
	str := fmt.Sprintf(BodyFormat, curTimeStr, remoteAddr,
		version.AppVer, version.GetShortGitCommit(), version.GetBuildTimeStr(), pkgutil.GetStartTime(), pkgddns.GetCurrentServiceIP())
	ctx.String(http.StatusOK, str)
}

func getFn1123CstTimeStr() string {
	t := pkgutil.GetCstNow()
	str := t.Format(Fn3339)
	return str
}
