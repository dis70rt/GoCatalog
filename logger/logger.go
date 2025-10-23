package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer Log.Sync()
}

func Error(msg string, err error) {
	if err != nil {
		Log.Error(msg, zap.Error(err))
	}
}

func Info(msg string) {
	Log.Info(msg)
}