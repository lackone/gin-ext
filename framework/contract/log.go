package contract

import (
	"context"
	"io"
	"time"
)

const LogKey = "ext:log"

type LogLevel uint8

const (
	UnknownLevel LogLevel = iota //未知日志级别
	PanicLevel                   //导致程序崩溃的日志
	FatalLevel                   //导致提前终止的错误
	ErrorLevel                   //错误，但不一定影响后续逻辑
	WarnLevel                    //警告，不影响后续逻辑
	InfoLevel                    //正常信息
	DebugLevel                   //调试信息
	TraceLevel                   //堆栈信息
)

// CtxFielder 定义了从context中获取信息的方法
type CtxFielder func(ctx context.Context) map[string]interface{}

// Formatter 定义了将日志信息组织成字符串的通用方法
type Formatter func(level LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error)

// Log 定义了日志服务协议
type Log interface {
	// Panic 表示会导致整个程序出现崩溃的日志信息
	Panic(ctx context.Context, msg string, fields map[string]interface{})
	// Fatal 表示会导致当前这个请求出现提前终止的错误信息
	Fatal(ctx context.Context, msg string, fields map[string]interface{})
	// Error 表示出现错误，但是不一定影响后续请求逻辑的错误信息
	Error(ctx context.Context, msg string, fields map[string]interface{})
	// Warn 表示出现错误，但是一定不影响后续请求逻辑的报警信息
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	// Info 表示正常的日志信息输出
	Info(ctx context.Context, msg string, fields map[string]interface{})
	// Debug 表示在调试状态下打印出来的日志信息
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	// Trace 表示最详细的信息，一般信息量比较大，可能包含调用堆栈等信息
	Trace(ctx context.Context, msg string, fields map[string]interface{})

	// SetLevel 设置日志级别
	SetLevel(level LogLevel)
	// SetCtxFielder 从context中获取上下文字段field
	SetCtxFielder(handler CtxFielder)
	// SetFormatter 设置输出格式
	SetFormatter(formatter Formatter)
	// SetOutput 设置输出管道
	SetOutput(out io.Writer)
}
