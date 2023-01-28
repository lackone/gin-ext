package distributed

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
)

type ExtDistributedLocalProvider struct {
}

func (e *ExtDistributedLocalProvider) Register(container framework.Container) framework.NewInstance {
	return NewExtDistributedLocal
}

func (e *ExtDistributedLocalProvider) Boot(container framework.Container) error {
	return nil
}

func (e *ExtDistributedLocalProvider) IsDefer() bool {
	return false
}

func (e *ExtDistributedLocalProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (e *ExtDistributedLocalProvider) Name() string {
	return contract.DistributedKey
}
