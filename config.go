package go_rss_reader

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"time"
)

var logger *zap.Logger

// InitZapLogger - init zap logger from JSON file and set as global zap
func InitZapLogger() (*zap.Logger, error) {
	jsonFile, errConfig := os.Open("zapConfig.json")
	if errConfig != nil {
		panic(errConfig)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cfg zap.Config
	if err := json.Unmarshal(byteValue, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	var err2 error
	logger, err2 = cfg.Build()
	if err2 != nil {
		panic(err2)
	}
	logger.Info("Zap logger construction succeeded")
	zap.ReplaceGlobals(logger)
	return logger, err2
}

// SyslogTimeEncoder - time stamp format for zap logger
func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan 01, 2006  15:04:05"))
}
