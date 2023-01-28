package services

import (
	"context"
	"github.com/lackone/gin-ext/framework"
	"github.com/lackone/gin-ext/framework/contract"
	"github.com/lackone/gin-ext/framework/provider/log/formatter"
	"io"
	pkgLog "log"
	"time"
)

type ExtLog struct {
	container  framework.Container //容器
	level      contract.LogLevel   //日志级别
	formatter  contract.Formatter  //格式化
	ctxFielder contract.CtxFielder //ctx上下文
	output     io.Writer           //输出
}

// IsLevelEnable 判断这个级别是否可以打印
func (e *ExtLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= e.level
}

// logf 为打印日志的核心函数
func (e *ExtLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	// 先判断日志级别
	if !e.IsLevelEnable(level) {
		return nil
	}

	// 使用ctxFielder 获取context中的信息
	fs := fields
	if e.ctxFielder != nil {
		t := e.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}

	// 如果绑定了trace服务，获取trace信息
	if e.container.IsBind(contract.TraceKey) {
		tracer := e.container.MustMake(contract.TraceKey).(contract.Trace)
		tc := tracer.GetTrace(ctx)
		if tc != nil {
			maps := tracer.ToMap(tc)
			for k, v := range maps {
				fs[k] = v
			}
		}
	}

	// 将日志信息按照formatter序列化为字符串
	if e.formatter == nil {
		e.formatter = formatter.TextFormatter
	}
	ct, err := e.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	// 如果是panic级别，则使用log进行panic
	if level == contract.PanicLevel {
		pkgLog.Panicln(string(ct))
		return nil
	}

	// 通过output进行输出
	e.output.Write(ct)
	e.output.Write([]byte("\r\n"))
	return nil
}

func (e *ExtLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.PanicLevel, ctx, msg, fields)
}

func (e *ExtLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.FatalLevel, ctx, msg, fields)
}

func (e *ExtLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.ErrorLevel, ctx, msg, fields)
}

func (e *ExtLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.WarnLevel, ctx, msg, fields)
}

func (e *ExtLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.InfoLevel, ctx, msg, fields)
}

func (e *ExtLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.DebugLevel, ctx, msg, fields)
}

func (e *ExtLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	e.logf(contract.TraceLevel, ctx, msg, fields)
}

func (e *ExtLog) SetLevel(level contract.LogLevel) {
	e.level = level
}

func (e *ExtLog) SetCtxFielder(handler contract.CtxFielder) {
	e.ctxFielder = handler
}

func (e *ExtLog) SetFormatter(formatter contract.Formatter) {
	e.formatter = formatter
}

func (e *ExtLog) SetOutput(out io.Writer) {
	e.output = out
}
