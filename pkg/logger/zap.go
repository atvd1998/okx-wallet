package logger

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var normalizeName = regexp.MustCompile(`[^A-Z0-9_]`).ReplaceAllString

// newZapConfig read log config from environment
// default value will be used if env is invalid
func newZapConfig(name string) *zap.Config {
	var (
		logLevel      = getPriorityEnv(name, "LOG_LEVEL", "debug")       // debug, info, warn, error, panic, fatal
		logColor      = getPriorityEnv(name, "LOG_COLOR", "f")           // true to enable
		logEncoding   = getPriorityEnv(name, "LOG_ENCODING", "console")  // console, json
		logTimestamp  = getPriorityEnv(name, "LOG_TIMESTAMP", "rfc3339") // rfc3339, rfc3339nano, iso8601, s, ms, ns, disabled
		logSeparator  = getPriorityEnv(name, "LOG_SEPARATOR", "\t")
		logStacktrace = getPriorityEnv(name, "LOG_STACKTRACE", "f")            // true to enable
		logSampling   = getPriorityEnv(name, "LOG_SAMPLING", "t")              // false to disable
		initial       = getPriorityEnv(name, "LOG_SAMPLING_INITIAL", "1000")   // at least 1000 per sec
		thereafter    = getPriorityEnv(name, "LOG_SAMPLING_THEREAFTER", "100") // 1 out of 100 per sec
	)

	// encoding
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.ConsoleSeparator = logSeparator
	encoderCfg.EncodeTime = parseTimeEncoding(logTimestamp)
	encoderCfg.TimeKey = parseTimeKey(logTimestamp)
	encoderCfg.EncodeLevel = parseLogEncoder(logColor)
	encoderCfg.NameKey = parseNameKey(name)

	// config
	config := zap.NewProductionConfig()
	config.Encoding = parseLogEncoding(logEncoding)
	config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))
	config.EncoderConfig = encoderCfg
	stacktrace, _ := strconv.ParseBool(logStacktrace)
	config.DisableStacktrace = !stacktrace
	config.Sampling = parseLogSampling(logSampling, initial, thereafter)
	return &config
}

// newSugaredLogger returns core logger as sugared logger with
// default basic config
func newSugaredLogger(name string) (*zap.SugaredLogger, error) {
	config := newZapConfig(name)
	log, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("error build logger / %w", err)
	}

	log = log.Named(name)
	return log.Sugar(), nil
}

func parseTimeEncoding(key string) zapcore.TimeEncoder {
	switch key {
	case "s":
		return unixTimeEncoder
	case "ms":
		return unixMilliTimeEncoder
	case "ns":
		return zapcore.EpochNanosTimeEncoder
	case "rfc3339nano":
		return zapcore.RFC3339NanoTimeEncoder
	case "rfc3339":
		return zapcore.RFC3339TimeEncoder
	case "iso8601":
		return zapcore.ISO8601TimeEncoder
	case "disabled":
		return func(time.Time, zapcore.PrimitiveArrayEncoder) {}
	default:
		return zapcore.RFC3339TimeEncoder
	}
}

func parseTimeKey(key string) string {
	if key == "disabled" {
		return zapcore.OmitKey
	}
	return "ts"
}

func parseLogEncoder(key string) zapcore.LevelEncoder {
	if enable, _ := strconv.ParseBool(key); enable {
		return zapcore.LowercaseColorLevelEncoder
	}
	return zapcore.LowercaseLevelEncoder
}

func parseLogSampling(key, initial, thereafter string) *zap.SamplingConfig {
	if enable, _ := strconv.ParseBool(key); !enable {
		return nil
	}
	return &zap.SamplingConfig{
		Initial:    toInt(initial, 1000),
		Thereafter: toInt(thereafter, 100),
	}
}

func toInt(val string, def int) int {
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}

func parseNameKey(key string) string {
	if key != "" {
		return "name"
	}
	return zapcore.OmitKey
}

func parseLogLevel(key string) zapcore.Level {
	switch strings.ToLower(key) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func parseLogEncoding(key string) string {
	if key == "json" {
		return key
	}
	return "console"
}

func unixTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

func unixMilliTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / int64(time.Millisecond))
}

// getEnv get key environment variable if exist otherwise return defalut Value
func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func getPriorityEnv(name, key, def string) string {
	if name == "" {
		return getEnv(key, def)
	}
	name = strings.ToUpper(name)
	name = normalizeName(name, "_")
	envKey := fmt.Sprintf("%s_%s", key, name)
	return getEnv(envKey, getEnv(key, def))
}
