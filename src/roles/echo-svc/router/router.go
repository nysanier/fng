package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nysanier/fng/src/roles/echo-svc/controller"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//v2 := r.Group("/v2")
	r.Use(gin.Recovery()) // 防止意外宕机
	r.GET("/", controller.Index)
	return r
}
