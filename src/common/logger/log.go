package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lowapple/elk/src/common/config"
	"log"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	S *zap.SugaredLogger // zap sugared logger
	L *zap.Logger        // zap logger
)

// Setup initialize the log instance
func Setup() bool {

	ok := true

	logLevel := parseLevel(config.Conf.Log.Level)

	cf := zap.NewProductionEncoderConfig()
	cf.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601-formatted (2022-05-06T17:14:21.101+0900) string with millisecond precision
	cEncoder := zapcore.NewConsoleEncoder(cf)
	var core zapcore.Core

	// 로그파일에 출력할 경우 파일과 콘솔출력
	if config.Conf.Log.WriteFile {
		fEncoder := zapcore.NewConsoleEncoder(cf) // NewJSONEncoder(cf)
		logFile := config.Conf.Log.Path + "/" + config.Conf.Log.FileName

		// set logrotate
		logf, err := rotatelogs.New(
			logFile+"."+config.Conf.Log.RotatePattern,
			rotatelogs.WithLinkName(logFile),                                                // 날짜가 없는 파일명으로 링크 생성
			rotatelogs.WithMaxAge(time.Duration(config.Conf.Log.RotateMaxAge)*24*time.Hour), // 보관일
			rotatelogs.WithRotationTime(time.Hour))                                          // 시간당 로테이트 동작
		if err != nil {
			log.Printf("error setup logging. %v", err)
			return false
		}

		fWriter := zapcore.AddSync(logf)

		core = zapcore.NewTee(
			zapcore.NewCore(fEncoder, fWriter, logLevel),
			zapcore.NewCore(cEncoder, zapcore.AddSync(os.Stdout), logLevel),
		)
	} else {
		// 콘솔만 출력하는 경우 (Docker 로 실행시 유용함)
		core = zapcore.NewTee(
			zapcore.NewCore(cEncoder, zapcore.AddSync(os.Stdout), logLevel),
		)
	}

	L = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	S = L.Sugar()
	return ok

}

// parseLevel string LogLevel 을 zapcore.Level 로 변환
//
// Parameters:
//   - lvl: log level (DEBUG, INFO, WARN, ERROR, FATAL)
//
// Return:
//   - Level
func parseLevel(lvl string) zapcore.Level {
	switch strings.ToUpper(lvl) {
	case "FATAL":
		return zapcore.FatalLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "WARN", "WARNING":
		return zapcore.WarnLevel
	case "INFO":
		return zapcore.InfoLevel
	case "DEBUG":
		return zapcore.DebugLevel
		//case "TRACE":
		//	return TRACE
	}

	log.Printf("not a valid log Level: %q", lvl)
	return zapcore.InfoLevel
}
