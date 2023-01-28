package config

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"path"
)

type ExtConfigProvider struct {
}

func (e *ExtConfigProvider) Register(container framework.Container) framework.NewInstance {
	return NewExtConfig
}

func (e *ExtConfigProvider) Boot(container framework.Container) error {
	return nil
}

func (e *ExtConfigProvider) IsDefer() bool {
	return false
}

func (e *ExtConfigProvider) Params(container framework.Container) []interface{} {
	app := container.MustMake(contract.AppKey).(contract.App)
	env := container.MustMake(contract.EnvKey).(contract.Env)

	appEnv := env.AppEnv()
	appEnvFolder := path.Join(app.ConfigFolder(), appEnv)

	return []interface{}{container, appEnvFolder, env.All()}
}

func (e *ExtConfigProvider) Name() string {
	return contract.ConfigKey
}
