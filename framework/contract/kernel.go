package contract

import "net/http"

const KernelKey = "ext:kernel"

type Kernel interface {
	//http.Handler结构，实际上是gin.Engine
	HttpEngine() http.Handler
}
