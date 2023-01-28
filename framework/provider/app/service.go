package app

import (
	"errors"
	"github.com/google/uuid"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/util"
	"path/filepath"
)

type ExtApp struct {
	container  framework.Container //服务容器
	baseFolder string              //基础目录
	appId      string              //表示当前app的唯一id
	configs    map[string]string   //app相关配置
}

func (e *ExtApp) AppId() string {
	return e.appId
}

func (e *ExtApp) Version() string {
	return "0.0.1"
}

func (e *ExtApp) BaseFolder() string {
	if e.baseFolder != "" {
		return e.baseFolder
	}

	return util.GetExecDirectory()
}

func (e *ExtApp) ConfigFolder() string {
	if val, ok := e.configs["config_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "config")
}

func (e *ExtApp) LogFolder() string {
	if val, ok := e.configs["log_folder"]; ok {
		return val
	}
	return filepath.Join(e.StorageFolder(), "log")
}

func (e *ExtApp) StorageFolder() string {
	if val, ok := e.configs["storage_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "storage")
}

func (e *ExtApp) ProviderFolder() string {
	if val, ok := e.configs["provider_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "app", "provider")
}

func (e *ExtApp) MiddlewareFolder() string {
	if val, ok := e.configs["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(e.HttpFolder(), "middleware")
}

func (e *ExtApp) CommandFolder() string {
	if val, ok := e.configs["command_folder"]; ok {
		return val
	}
	return filepath.Join(e.ConsoleFolder(), "command")
}

func (e *ExtApp) RuntimeFolder() string {
	if val, ok := e.configs["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(e.StorageFolder(), "runtime")
}

func (e *ExtApp) TestFolder() string {
	if val, ok := e.configs["test_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "test")
}

func (e *ExtApp) HttpFolder() string {
	if val, ok := e.configs["http_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "app", "http")
}

func (e *ExtApp) ConsoleFolder() string {
	if val, ok := e.configs["console_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "app", "console")
}

func (e *ExtApp) AppFolder() string {
	if val, ok := e.configs["app_folder"]; ok {
		return val
	}
	return filepath.Join(e.BaseFolder(), "app")
}

func (e *ExtApp) LoadAppConfig(configs map[string]string) {
	for k, v := range configs {
		e.configs[k] = v
	}
}

func newExtApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	container := params[0].(framework.Container)
	baseFolder := params[1].(string)

	appId := uuid.New().String()

	return &ExtApp{container: container, baseFolder: baseFolder, appId: appId}, nil
}
