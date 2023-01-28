package id

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
)

type ExtIdProvider struct {
}

func (e *ExtIdProvider) Register(container framework.Container) framework.NewInstance {
	return NewExtId
}

func (e *ExtIdProvider) Boot(container framework.Container) error {
	return nil
}

func (e *ExtIdProvider) IsDefer() bool {
	return false
}

func (e *ExtIdProvider) Params(container framework.Container) []interface{} {
	return []interface{}{}
}

func (e *ExtIdProvider) Name() string {
	return contract.IdKey
}
