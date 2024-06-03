package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func InitLogger() {

	// 创建并配置文件输出
	logFileName := "tmp/logs/" + time.Now().Format("2006-01-02T15-04-05") + ".log"
	file, err := os.Create(logFileName)
	if err != nil {
		panic(err)
	}

	// 定义日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	// 配置 Console 和 File 输出
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	// 创建 zapcore.Core 每个输出路径
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), atomicLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), atomicLevel),
	)

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.Development())

	zap.ReplaceGlobals(logger)
}
