package main

import (
	"github.com/lackone/gin-ext/app/console"
	"github.com/lackone/gin-ext/app/http"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/provider/app"
	"github.com/lackone/gin-ext/framework/provider/config"
	"github.com/lackone/gin-ext/framework/provider/distributed"
	"github.com/lackone/gin-ext/framework/provider/env"
	"github.com/lackone/gin-ext/framework/provider/id"
	"github.com/lackone/gin-ext/framework/provider/kernel"
	"github.com/lackone/gin-ext/framework/provider/log"
	"github.com/lackone/gin-ext/framework/provider/trace"
)

func main() {
	//初始化服务容器
	container := framework.NewExtContainer()

	//绑定APP服务提代者
	container.Bind(&app.ExtAppProvider{})
	//后续初始化需要绑定的服务提供者
	container.Bind(&env.ExtEnvProvider{})
	container.Bind(&config.ExtConfigProvider{})
	container.Bind(&log.ExtLogProvider{})
	container.Bind(&id.ExtIdProvider{})
	container.Bind(&trace.ExtTraceProvider{})
	container.Bind(&distributed.ExtDistributedLocalProvider{})

	//将HTTP引擎，绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.ExtKernelProvider{HttpEngine: engine})
	}

	//运行命令
	console.RunCommand(container)
}
