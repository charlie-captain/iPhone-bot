package log

import (
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				//处理文件名
				fileName := path.Base(frame.File)
				return frame.Function, fileName
			},
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.DebugLevel,
	}
	Log.SetReportCaller(true)
	Log.Println("init log")
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Log.Fatal(err)
	}
	// defer file.Close()
	// 创建新的log对象
	writers := []io.Writer{
		file,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	Log.SetOutput(fileAndStdoutWriter)
}
