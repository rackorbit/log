package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

// Config represents the configuration for a Logger.
type Config struct {
	// Encoder is the encoder that should be used to encode log output.
	Encoder zapcore.Encoder `json:"encoder"`
	// Level is the level at which logs will be output.
	Level Level `json:"level"`
	// Output is the writer that encoded logs will be output to.
	Output io.Writer `json:"output"`
}
