package klog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

type Klog struct {
	//logger *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

//func (kl *Klog)StartUp()  {
//	kl.InitLogger()
//	defer kl.logger.Sync()
//	kl.simpleHttpGet("www.google.com")
//	kl.simpleHttpGet("http://www.google.com")
//}
//
//func (kl *Klog)InitLogger() {
//	kl.logger, _ = zap.NewProduction()
//}
//
//func (kl *Klog)simpleHttpGet(url string) {
//	resp, err := http.Get(url)
//	if err != nil {
//		kl.logger.Error(
//			"Error fetching url..",
//			zap.String("url", url),
//			zap.Error(err))
//	} else {
//		kl.logger.Info("Success..",
//			zap.String("statusCode", resp.Status),
//			zap.String("url", url))
//		resp.Body.Close()
//	}
//}





func (kl *Klog)StartUpSugared()  {
	kl.InitLoggerSugared()
	defer kl.sugaredLogger.Sync()
	kl.simpleSugaredHttpGet("www.google.com")
	kl.simpleSugaredHttpGet("http://www.google.com")
}

func (kl *Klog)InitLoggerSugared() {

	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	curLogger := zap.New(core)
	kl.sugaredLogger = curLogger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}



func (kl *Klog)simpleSugaredHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		kl.sugaredLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		kl.sugaredLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}
