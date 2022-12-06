package log

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"unsafe"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger unsafe.Pointer
)

func L() *Logger {
	l := (*Logger)(atomic.LoadPointer(&defaultLogger))
	if l == nil {
		l = NewLogger()
		atomic.StorePointer(&defaultLogger, unsafe.Pointer(l))
	}
	return l
}

type Logger struct {
	base  *zap.Logger
	level zap.AtomicLevel
}

func NewLogger() *Logger {
	logger := &Logger{level: zap.NewAtomicLevelAt(zapcore.InfoLevel)}
	logger.base = zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(productionEncoderConfig),
			nopCloserSink{os.Stderr},
			logger.level,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	return logger
}

func NewDevelopmentLogger() *Logger {
	logger := &Logger{level: zap.NewAtomicLevelAt(zapcore.DebugLevel)}
	logger.base = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(developmentEncoderConfig),
			nopCloserSink{os.Stderr},
			logger.level,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.WarnLevel),
		zap.Development(),
	)
	return logger
}

func NewNopLogger() *Logger {
	return &Logger{zap.NewNop(), zap.NewAtomicLevel()}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...Field) *Logger {
	l.base = l.base.With(fields...)
	return l
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.base.Sync()
}

// Level returns the minimum enabled log level.
func (l *Logger) Level() Level {
	return l.level.Level()
}

// SetLevel alters the logging level.
func (l *Logger) SetLevel(level Level) {
	l.level.SetLevel(level)
}

// Named adds a new path segment to the Logger's name and return the new Logger.
// Segments are joined by periods. By default, Logger are unnamed.
func (l *Logger) Named(name string) *Logger {
	if name == `` {
		return l
	}
	c := *l
	c.base = l.base.Named(name)
	return &c
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.DebugLevel, msg, kv)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.InfoLevel, msg, kv)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.WarnLevel, msg, kv)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.ErrorLevel, msg, kv)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See zapcore.DPanicLevel for details.)
func (l *Logger) DPanic(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.DPanicLevel, msg, kv)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.PanicLevel, msg, kv)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(ctx context.Context, msg string, kv ...interface{}) {
	kv = append(kv, TraceId(ctx))
	l.logw(zapcore.FatalLevel, msg, kv)
}

//Deprecated: Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(zapcore.DebugLevel, format, args)
}

//Deprecated: Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(zapcore.InfoLevel, format, args)
}

//Deprecated: Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(zapcore.WarnLevel, format, args)
}

//Deprecated: Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(zapcore.ErrorLevel, format, args)
}

//Deprecated: DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.logf(zapcore.DPanicLevel, format, args)
}

//Deprecated: Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logf(zapcore.PanicLevel, format, args)
}

//Deprecated: Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(zapcore.FatalLevel, format, args)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (l *Logger) Debugw(msg string, kv ...interface{}) {
	l.logw(zapcore.DebugLevel, msg, kv)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Infow(msg string, kv ...interface{}) {
	l.logw(zapcore.InfoLevel, msg, kv)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Warnw(msg string, kv ...interface{}) {
	l.logw(zapcore.WarnLevel, msg, kv)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Errorw(msg string, kv ...interface{}) {
	l.logw(zapcore.ErrorLevel, msg, kv)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See zapcore.DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) DPanicw(msg string, kv ...interface{}) {
	l.logw(zapcore.InfoLevel, msg, kv)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *Logger) Panicw(msg string, kv ...interface{}) {
	l.logw(zapcore.PanicLevel, msg, kv)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *Logger) Fatalw(msg string, kv ...interface{}) {
	l.logw(zapcore.FatalLevel, msg, kv)
}

func (l *Logger) logf(lvl zapcore.Level, format string, args []interface{}) {
	if lvl < zapcore.DPanicLevel && !l.base.Core().Enabled(lvl) {
		return
	}
	if len(args) > 0 {
		if format == `` {
			format = fmt.Sprint(args...)
		} else {
			format = fmt.Sprintf(format, args...)
		}
	}
	if ce := l.base.Check(lvl, format); ce != nil {
		ce.Write()
	}
}

func (l *Logger) logw(lvl zapcore.Level, msg string, kv []interface{}) {
	if lvl < zapcore.DPanicLevel && !l.base.Core().Enabled(lvl) {
		return
	}
	if ce := l.base.Check(lvl, msg); ce != nil {
		if n := len(kv); n > 0 {
			fields, invalids := make([]zapcore.Field, 0, n), invalidPairs(nil)

			for i, m := 0, n-1; i < n; {
				if f, ok := kv[i].(zapcore.Field); ok {
					fields = append(fields, f)
					i++
					continue
				}

				if ctx, ok := kv[i].(context.Context); ok {
					i++

					f := TraceId(ctx)
					if f.String != NoTraceId {
						fields = append(fields, TraceId(ctx))
					}

					continue
				}

				if i == m {
					l.base.DPanic(danglingKeyErrMsg, zap.Any(`ignored`, kv[i]))
					break
				}
				k, v := kv[i], kv[i+1]
				if s, ok := k.(string); !ok {
					if cap(invalids) == 0 {
						invalids = make(invalidPairs, 0, n/2)
					}
					invalids = append(invalids, invalidPair{i, k, v})
				} else {
					fields = append(fields, zap.Any(s, v))
				}
				i += 2
			}
			if len(invalids) > 0 {
				l.base.DPanic(nonStringKeyErrMsg, zap.Array(`invalid`, invalids))
			}
			ce.Write(fields...)
		} else {
			ce.Write()
		}
	}
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64(`position`, int64(p.position))
	zap.Any(`key`, p.key).AddTo(enc)
	zap.Any(`value`, p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}

type nopCloserSink struct{ zapcore.WriteSyncer }

func (nopCloserSink) Close() error { return nil }

const (
	danglingKeyErrMsg  = `Ignored key without a value.`
	nonStringKeyErrMsg = `Ignored key-value pairs with non-string keys.`
)
