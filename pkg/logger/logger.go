package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Entry struct {
	logger *zap.SugaredLogger
	fields map[string]interface{}
}

// 全局日志实例和配置
var (
	// defaultLogger 是默认的日志记录器
	defaultLogger *zap.Logger
	// defaultSugar 是默认的SugaredLogger
	defaultSugar *zap.SugaredLogger
	// defaultLevel 是默认的日志级别
	defaultLevel = zapcore.InfoLevel
)

// LogLevel 日志级别定义
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// 将字符串日志级别转换为zapcore级别
func parseLevel(level LogLevel) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "component",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色日志级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短的调用者路径
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		defaultLevel,
	)

	// 创建Logger
	defaultLogger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1), // 跳过封装函数调用
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	defaultSugar = defaultLogger.Sugar()

	fmt.Println("日志系统初始化完成，当前日志级别:", defaultLevel)
}

func SetLevel(level LogLevel) {
	defaultLevel = parseLevel(level)
	// 重建Logger以应用新级别
	Init()
}

func WithField(key string, value interface{}) *Entry {
	return &Entry{
		logger: defaultSugar,
		fields: map[string]interface{}{key: value},
	}
}

func WithFields(fields map[string]interface{}) *Entry {
	return &Entry{
		logger: defaultSugar,
		fields: fields,
	}
}

// 私有方法，用于实际记录日志
func (e *Entry) log(level zapcore.Level, msg string, args ...interface{}) {
	// 如果没有额外字段，直接使用默认logger
	if len(e.fields) == 0 {
		// 直接使用全局logger
		switch level {
		case zapcore.DebugLevel:
			defaultSugar.Debugf(msg, args...)
		case zapcore.InfoLevel:
			defaultSugar.Infof(msg, args...)
		case zapcore.WarnLevel:
			defaultSugar.Warnf(msg, args...)
		case zapcore.ErrorLevel:
			defaultSugar.Errorf(msg, args...)
		case zapcore.FatalLevel:
			defaultSugar.Fatalf(msg, args...)
		}
		return
	}

	// 创建带有额外字段的logger
	fields := make([]interface{}, 0, len(e.fields)*2)
	for k, v := range e.fields {
		fields = append(fields, k, v)
	}

	with := defaultSugar.With(fields...)

	switch level {
	case zapcore.DebugLevel:
		with.Debugf(msg, args...)
	case zapcore.InfoLevel:
		with.Infof(msg, args...)
	case zapcore.WarnLevel:
		with.Warnf(msg, args...)
	case zapcore.ErrorLevel:
		with.Errorf(msg, args...)
	case zapcore.FatalLevel:
		with.Fatalf(msg, args...)
	}
}

func (e *Entry) Debug(msg string, args ...interface{}) {
	e.log(zapcore.DebugLevel, msg, args...)
}

func (e *Entry) Info(msg string, args ...interface{}) {
	e.log(zapcore.InfoLevel, msg, args...)
}

func (e *Entry) Warn(msg string, args ...interface{}) {
	e.log(zapcore.WarnLevel, msg, args...)
}

func (e *Entry) Error(msg string, args ...interface{}) {
	e.log(zapcore.ErrorLevel, msg, args...)
}

func (e *Entry) Fatal(msg string, args ...interface{}) {
	e.log(zapcore.FatalLevel, msg, args...)
}

// Info 封装的信息级别日志
func Info(msg string, args ...interface{}) {
	green := "\033[32m"
	reset := "\033[0m"
	defaultSugar.Infof(green+msg+reset, args...)
}

// Debug 封装的调试级别日志
func Debug(msg string, args ...interface{}) {
	defaultSugar.Debugf(msg, args...)
}

// Warn 封装的警告级别日志
func Warn(msg string, args ...interface{}) {
	defaultSugar.Warnf(msg, args...)
}

// Error 封装的错误级别日志
func Error(msg string, args ...interface{}) {
	defaultSugar.Errorf(msg, args...)
}

// Fatal 封装的致命错误日志，记录后会导致程序退出
func Fatal(msg string, args ...interface{}) {
	defaultSugar.Fatalf(msg, args...)
}
