package log

import (
	"context"
	"jonghong/internal/pkg/known"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugw(msg string, keysAndValues ...any) //Debug级别
	Infow(msg string, keysAndValues ...any)  // Info几笔
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Panicw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
	Sync()
}

type zapLogger struct {
	z *zap.Logger
}

// 检查接口是否实现
var _ Logger = &zapLogger{}

func (zl *zapLogger) Debugw(msg string, keysAndValues ...any) {
	zl.z.Sugar().Debugw(msg, keysAndValues...)
}

func (zl *zapLogger) Infow(msg string, keysAndValues ...any) {
	zl.z.Sugar().Infow(msg, keysAndValues...)
}

func (zl *zapLogger) Warnw(msg string, keysAndValues ...any) {
	zl.z.Sugar().Warnw(msg, keysAndValues...)
}

func (zl *zapLogger) Errorw(msg string, keysAndValues ...any) {
	zl.z.Sugar().Errorw(msg, keysAndValues...)
}

func (zl *zapLogger) Panicw(msg string, keysAndValues ...any) {
	zl.z.Sugar().Panicw(msg, keysAndValues...)
}

func (zl *zapLogger) Fatalw(msg string, keysAndValues ...any) {
	zl.z.Sugar().Fatalw(msg, keysAndValues...)
}

func (zl *zapLogger) Sync() {
	zl.z.Sync()
}

var (
	mu  sync.Mutex                // 用于上锁，防止std对象的重复创建
	std = NewLogger(NewOptions()) // 按照默认配置初始化一个logger对象
)

func Init(opt *Options) {
	mu.Lock()
	defer mu.Unlock() // 代码块运行结束后执行
	std = NewLogger(opt)
}

// 初始化的方法
func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 解析日志的等级
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	// 创建新的编码器：用于转换日期变量
	encoderConfig := zap.NewProductionEncoderConfig()

	// 暂时不清楚有啥用
	encoderConfig.MessageKey = "message"
	// 暂时不清楚有啥用
	encoderConfig.TimeKey = "timestamp"
	// 格式化日期 yyyy-mm-dd hh:MM:ss
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2025-01-01 00:00:00.000"))
	}
	// 格式化持续时间为毫秒
	encoderConfig.EncodeDuration = func(d time.Duration, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendFloat64(float64(d) / float64(time.Millisecond))
	}
	// 设置日志的配置项
	cfg := &zap.Config{
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
	}
	// 创建zap对象
	z, err := cfg.Build(
		// 允许panic等级以上的日志输出堆栈信息
		zap.AddStacktrace(zapcore.PanicLevel),
		// 日志打印时显示日志函数调用者的位置
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{z: z}
	// 把标准库的 log.Logger 的 info 级别的输出重定向到 zap.Logger
	zap.RedirectStdLog(z)

	return logger

}

// 定义包层面的方法

func Debugw(msg string, keysAndValues ...any) {
	std.z.Sugar().Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...any) {
	std.z.Sugar().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...any) {
	std.z.Sugar().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...any) {
	std.z.Sugar().Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...any) {
	std.z.Sugar().Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...any) {
	std.z.Sugar().Fatalw(msg, keysAndValues...)
}

func Sync() {
	std.z.Sync()
}

func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}

func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()
	requestID := ctx.Value(known.XRequestIDKey)
	if requestID != nil {
		lc.z = lc.z.With(zap.Any(known.XRequestIDKey, requestID))
	}
	if userID := ctx.Value(known.XUsernameKey); userID != nil {
		lc.z = lc.z.With(zap.Any(known.XUsernameKey, userID))
	}
	return lc
}

func (l *zapLogger) clone() *zapLogger {
	lc := *l
	return &lc
}
