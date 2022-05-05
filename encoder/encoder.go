// Package encoder implements log encoders used to output logs in varying ways.
package encoder

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Encoder is an interface for a log encoder.
type Encoder interface {
	zapcore.Encoder
	Provision() error
}

// Config represents the configuration for a log encoder.
type Config struct {
	MessageKey     *string `json:"message_key,omitempty"`
	LevelKey       *string `json:"level_key,omitempty"`
	TimeKey        *string `json:"time_key,omitempty"`
	NameKey        *string `json:"name_key,omitempty"`
	CallerKey      *string `json:"caller_key,omitempty"`
	StacktraceKey  *string `json:"stacktrace_key,omitempty"`
	LineEnding     *string `json:"line_ending,omitempty"`
	TimeFormat     string  `json:"time_format,omitempty"`
	DurationFormat string  `json:"duration_format,omitempty"`
	LevelFormat    string  `json:"level_format,omitempty"`
	CallerFormat   string  `json:"caller_format,omitempty"`
}

// ZapcoreEncoderConfig returns a zap encoder config, built from the Config
// values.
func (c *Config) ZapcoreEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()

	if c == nil {
		c = new(Config)
	}
	if c.MessageKey != nil {
		cfg.MessageKey = *c.MessageKey
	}
	if c.LevelKey != nil {
		cfg.LevelKey = *c.LevelKey
	}
	if c.TimeKey != nil {
		cfg.TimeKey = *c.TimeKey
	}
	if c.NameKey != nil {
		cfg.NameKey = *c.NameKey
	}
	if c.CallerKey != nil {
		cfg.CallerKey = *c.CallerKey
	}
	if c.StacktraceKey != nil {
		cfg.StacktraceKey = *c.StacktraceKey
	}
	if c.LineEnding != nil {
		cfg.LineEnding = *c.LineEnding
	}

	var timeFormatter zapcore.TimeEncoder
	switch c.TimeFormat {
	case "":
		timeFormatter = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("15:04:05"))
		}
	case "unix_seconds_float":
		timeFormatter = zapcore.EpochTimeEncoder
	case "unix_milli_float":
		timeFormatter = zapcore.EpochMillisTimeEncoder
	case "unix_nano":
		timeFormatter = zapcore.EpochNanosTimeEncoder
	case "iso8601":
		timeFormatter = zapcore.ISO8601TimeEncoder
	default:
		timeFormat := c.TimeFormat
		switch c.TimeFormat {
		case "rfc3339":
			timeFormat = time.RFC3339
		case "rfc3339_nano":
			timeFormat = time.RFC3339Nano
		case "wall":
			timeFormat = "2006/01/02 15:04:05"
		case "wall_milli":
			timeFormat = "2006/01/02 15:04:05.000"
		case "wall_nano":
			timeFormat = "2006/01/02 15:04:05.000000000"
		case "common_log":
			timeFormat = "02/Jan/2006:15:04:05 -0700"
		}
		timeFormatter = func(ts time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(ts.UTC().Format(timeFormat))
		}
	}
	cfg.EncodeTime = timeFormatter

	var durFormatter zapcore.DurationEncoder
	switch c.DurationFormat {
	case "seconds":
		durFormatter = zapcore.SecondsDurationEncoder
	case "nano":
		durFormatter = zapcore.NanosDurationEncoder
	case "", "string":
		durFormatter = zapcore.StringDurationEncoder
	}
	cfg.EncodeDuration = durFormatter

	var levelFormatter zapcore.LevelEncoder
	switch c.LevelFormat {
	case "lower":
		levelFormatter = zapcore.LowercaseLevelEncoder
	case "upper":
		levelFormatter = zapcore.CapitalLevelEncoder
	case "", "color":
		levelFormatter = zapcore.CapitalColorLevelEncoder
	}
	cfg.EncodeLevel = levelFormatter

	var callerFormatter zapcore.CallerEncoder
	switch c.CallerFormat {
	case "short":
		callerFormatter = zapcore.ShortCallerEncoder
	case "full":
		callerFormatter = zapcore.FullCallerEncoder
	}
	cfg.EncodeCaller = callerFormatter

	return cfg
}
