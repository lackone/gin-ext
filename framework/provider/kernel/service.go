package kernel

import (
	"github.com/lackone/gin-ext/framework/gin"
	"net/http"
) /**/

type ExtKernel struct {
	engine *gin.Engine
}

func NewExtKernel(params ...interface{}) (interface{}, error) {
	engine := params[0].(*gin.Engine)
	return &ExtKernel{engine: engine}, nil
}

func (e *ExtKernel) HttpEngine() http.Handler {
	return e.engine
}
