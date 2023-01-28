package env

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
)

type ExtEnvProvider struct {
	folder string //.env所在目录
}

func (e *ExtEnvProvider) Register(container framework.Container) framework.NewInstance {
	return NewExtEnv
}

func (e *ExtEnvProvider) Boot(container framework.Container) error {
	app := container.MustMake(contract.AppKey).(contract.App)
	e.folder = app.BaseFolder()
	return nil
}

func (e *ExtEnvProvider) IsDefer() bool {
	return false
}

func (e *ExtEnvProvider) Params(container framework.Container) []interface{} {
	return []interface{}{e.folder}
}

func (e *ExtEnvProvider) Name() string {
	return contract.EnvKey
}
