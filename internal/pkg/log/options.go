package log

import "go.uber.org/zap/zapcore"

type Options struct {
	// 是否打印行号
	DisableCaller bool
	// 是否打印堆栈信息
	DisableStacktrace bool
	// 默认的日志界别
	Level string
	// 日志的输出方式，console，log
	Format string
	// 日志输出地址 file, stdout
	OutputPaths []string
}

func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             zapcore.InfoLevel.String(),
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}
