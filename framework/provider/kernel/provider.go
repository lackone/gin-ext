package kernel

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/gin"
)

type ExtKernelProvider struct {
	HttpEngine *gin.Engine
}

func (e *ExtKernelProvider) Register(container framework.Container) framework.NewInstance {
	return NewExtKernel
}

func (e *ExtKernelProvider) Boot(container framework.Container) error {
	if e.HttpEngine == nil {
		e.HttpEngine = gin.Default()
	}
	e.HttpEngine.SetContainer(container)
	return nil
}

func (e *ExtKernelProvider) IsDefer() bool {
	return false
}

func (e *ExtKernelProvider) Params(container framework.Container) []interface{} {
	return []interface{}{e.HttpEngine}
}

func (e *ExtKernelProvider) Name() string {
	return contract.KernelKey
}
