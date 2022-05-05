package log

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/rackorbit/log/encoder"
)

// Flags are CLI flags for use with [kong](https://github.com/alecthomas/kong)
type Flags struct {
	Encoder string `enum:"console,json,logfmt" default:"logfmt"`
	Level   string `enum:"debug,info,warn,error,panic,fatal" default:"info"`
}

// RegisterFromFlags creates and registers a new logger using CLI flags.
func RegisterFromFlags(f Flags) (func(), error) {
	level, ok := Levels[f.Level]
	if !ok {
		return nil, fmt.Errorf("log: unknown log level: \"%s\"", f.Level)
	}

	var e encoder.Encoder
	switch f.Encoder {
	case "console":
		e = &encoder.Console{}
	case "json":
		e = &encoder.JSON{}
	case "logfmt":
		e = &encoder.Logfmt{}
	default:
		return nil, fmt.Errorf("log: unknown log encoder: \"%s\"", f.Encoder)
	}

	if err := e.Provision(); err != nil {
		return nil, errors.Wrap(err, "log: failed to provision encoder")
	}

	l, err := New(&Config{
		Encoder: e,
		Level:   level,
		Output:  os.Stderr,
	})
	if err != nil {
		return nil, errors.Wrap(err, "log: failed to configure logger")
	}
	if err := Register(l); err != nil {
		return nil, errors.Wrap(err, "log: failed to register logger")
	}
	return Sync, nil
}
