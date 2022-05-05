package encoder

import (
	"go.uber.org/zap/zapcore"
)

// Console is a log encoder for human-readable colored console output. This
// encoder is designed for development where only humans are reading the output.
type Console struct {
	Config
	zapcore.Encoder `json:"-"`
}

var _ Encoder = (*Console)(nil)

func (e *Console) Provision() error {
	e.CallerFormat = "short"

	e.Encoder = zapcore.NewConsoleEncoder(e.Config.ZapcoreEncoderConfig())
	return nil
}
