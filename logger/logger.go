package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化日志配置
func InitLogger(level string, logPath string) error {
	// 创建日志目录
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 设置日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置日志输出和轮转
	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // 每个日志文件最大尺寸（MB）
		MaxBackups: 3,   // 保留的旧日志文件最大数量
		MaxAge:     7,   // 保留的旧日志文件最大天数
		Compress:   true,
	}

	// 设置日志编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zapLevel,
	)

	// 创建logger
	Logger = zap.New(core, zap.AddCaller())
	defer Logger.Sync()

	Logger.Info("日志系统初始化成功",
		zap.String("level", level),
		zap.String("path", logPath),
		zap.Time("time", time.Now()),
	)

	return nil
}
