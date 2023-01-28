package http

import "github.com/lackone/gin-ext/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	//业务路由操作
	Routes(engine)
	return engine, nil
}
