package trace

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
)

type ExtTraceProvider struct {
	container framework.Container
}

func (e *ExtTraceProvider) Register(container framework.Container) framework.NewInstance {
	return NewExtTrace
}

func (e *ExtTraceProvider) Boot(container framework.Container) error {
	e.container = container
	return nil
}

func (e *ExtTraceProvider) IsDefer() bool {
	return false
}

func (e *ExtTraceProvider) Params(container framework.Container) []interface{} {
	return []interface{}{e.container}
}

func (e *ExtTraceProvider) Name() string {
	return contract.TraceKey
}
