package vch

import (
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var MyLog *logrus.Entry

func init() {
	MyLog = logrus.NewEntry(logrus.StandardLogger())
}

func InitLog(path string) error {
	logPath := path // 日志存放路径
	logPath += `/`
	mqttLogger := logrus.New()
	mqttLogger.Out = io.Discard
	mqttLogger.SetLevel(logrus.DebugLevel)
	// 获取保存周期，最小为7天
	maxAge := 7 * 24 * time.Hour
	if maxAge < 7*24*time.Hour {
		maxAge = 7 * 24 * time.Hour
	}

	// 获取切割周期，最小为一个小时
	rotatimeTime := time.Hour
	if rotatimeTime < time.Hour {
		rotatimeTime = time.Hour
	}
	debugWriter, _ := rotatelogs.New(
		logPath+"/vch.log",                        // 日志文件名格式
		rotatelogs.WithMaxAge(maxAge),             // 设置最大保存周期
		rotatelogs.WithRotationTime(rotatimeTime), // 设置日志切割周期，最小为1小时
	)
	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: debugWriter,
		logrus.FatalLevel: debugWriter,
		logrus.ErrorLevel: debugWriter,
		logrus.WarnLevel:  debugWriter,
		logrus.InfoLevel:  debugWriter, // info级别使用logWriter写日志
		logrus.DebugLevel: debugWriter,
		logrus.TraceLevel: debugWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano, // 格式日志时间
	})
	mqttLogger.AddHook(Hook)
	MyLog = logrus.NewEntry(mqttLogger)
	return nil
}
