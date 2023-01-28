package app

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
)

type ExtAppProvider struct {
	BaseFolder string
}

func (e *ExtAppProvider) Register(container framework.Container) framework.NewInstance {
	return newExtApp
}

func (e *ExtAppProvider) Boot(container framework.Container) error {
	return nil
}

func (e *ExtAppProvider) IsDefer() bool {
	return false
}

func (e *ExtAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, e.BaseFolder}
}

func (e *ExtAppProvider) Name() string {
	return contract.AppKey
}
