package log

import (
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/provider/log/formatter"
	"github.com/lackone/gin-ext/framework/provider/log/services"
	"io"
	"strings"
)

type ExtLogProvider struct {
	Driver string // Driver

	// 日志级别
	Level contract.LogLevel
	// 日志输出格式方法
	Formatter contract.Formatter
	// 日志context上下文信息获取函数
	CtxFielder contract.CtxFielder
	// 日志输出信息
	Output io.Writer
}

func (e *ExtLogProvider) Register(container framework.Container) framework.NewInstance {
	if e.Driver == "" {
		conf, err := container.Make(contract.ConfigKey)
		if err != nil {
			// 默认使用console
			return services.NewExtConsoleLog
		}

		cs := conf.(contract.Config)
		e.Driver = strings.ToLower(cs.GetString("log.driver"))
	}

	// 根据driver的配置项确定
	switch e.Driver {
	case "single":
		return services.NewExtSingleLog
	case "rotate":
		return services.NewExtRotateLog
	case "console":
		return services.NewExtConsoleLog
	case "custom":
		return services.NewExtCustomLog
	default:
		return services.NewExtConsoleLog
	}
}

func (e *ExtLogProvider) Boot(container framework.Container) error {
	return nil
}

func (e *ExtLogProvider) IsDefer() bool {
	return false
}

func (e *ExtLogProvider) Params(container framework.Container) []interface{} {
	// 获取configService
	conf := container.MustMake(contract.ConfigKey).(contract.Config)

	// 设置参数formatter
	if e.Formatter == nil {
		e.Formatter = formatter.TextFormatter
		if conf.IsExist("log.formatter") {
			v := conf.GetString("log.formatter")
			if v == "json" {
				e.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				e.Formatter = formatter.TextFormatter
			}
		}
	}

	if e.Level == contract.UnknownLevel {
		e.Level = contract.InfoLevel
		if conf.IsExist("log.level") {
			e.Level = logLevel(conf.GetString("log.level"))
		}
	}

	// 定义5个参数
	return []interface{}{container, e.Level, e.CtxFielder, e.Formatter, e.Output}
}

func (e *ExtLogProvider) Name() string {
	return contract.LogKey
}

// logLevel get level from string
func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
