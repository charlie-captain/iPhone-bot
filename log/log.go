package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
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
