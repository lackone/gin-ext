package http

import (
    "github.com/lackone/gin-ext/app/http/modules/demo"
    "github.com/lackone/gin-ext/framework/gin"
    "github.com/lackone/gin-ext/framework/middleware/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {
    // /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
    r.Use(static.Serve("/", static.LocalFile("./dist", false)))

    demo.Register(r)
}
