// Package log is a wrapper around the zap logger, providing easier integration
// of logging into applications and packages.
package log

import (
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_mx     sync.Mutex
	_once   sync.Once
	_logger *zap.Logger
)

// Register registers a Logger globally. Register can be called multiple times,
// but only the first call will register a logger, all subsequent calls are
// a no-op.
func Register(l *Logger) error {
	_mx.Lock()
	defer _mx.Unlock()

	p, err := l.Logger()
	if err != nil {
		return errors.Wrap(err, "log: failed to get zap logger")
	}

	_once.Do(func() {
		_logger = p
		zap.RedirectStdLog(p)
	})
	return nil
}

// Logger likes logs.
type Logger struct {
	config *Config
}

// New returns a new configured Logger.
func New(c *Config) (*Logger, error) {
	l := &Logger{
		config: c,
	}
	if l.config.Encoder == nil {
		return nil, errors.New("log: no encoder configured")
	}
	if l.config.Output == nil {
		return nil, errors.New("log: no output configured")
	}
	return l, nil
}

// Logger returns a zap Logger that can be used to actually log things.
func (l *Logger) Logger() (*zap.Logger, error) {
	return zap.New(
		zapcore.NewCore(
			l.config.Encoder,
			zapcore.AddSync(l.config.Output),
			l.config.Level,
		),
		// Add the caller field to log output.
		zap.AddCaller(),
		// Required because logger function calls are routed through this package.
		zap.AddCallerSkip(1),
		// Only log stacktraces on panics or higher, default is for warnings.
		zap.AddStacktrace(zap.PanicLevel),
	), nil
}

// Log returns the globally configured logger instance.
func Log() *zap.Logger {
	return _logger
}

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return _logger.Check(lvl, msg)
}

// Core returns the Logger's underlying zapcore.Core.
func Core() zapcore.Core {
	return _logger.Core()
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func Named(s string) *zap.Logger {
	return _logger.Named(s)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() {
	_ = _logger.Sync()
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return _logger.WithOptions(zap.AddCallerSkip(-1)).With(fields...)
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func WithOptions(opts ...zap.Option) *zap.Logger {
	return _logger.WithOptions(opts...)
}
