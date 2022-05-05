package encoder

import (
	"go.uber.org/zap/zapcore"
)

// JSON is a log encoder for machine-readable JSON objects. Each line is output
// as a separate JSON object, essentially making the entire log output JSON-ND.
// This encoder is designed for production where logs are read and ingested by
// log aggregation systems.
type JSON struct {
	Config
	zapcore.Encoder `json:"-"`
}

var _ Encoder = (*JSON)(nil)

func (e *JSON) Provision() error {
	e.CallerFormat = "full"

	e.Encoder = zapcore.NewJSONEncoder(e.Config.ZapcoreEncoderConfig())
	return nil
}
