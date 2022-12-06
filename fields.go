package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// A Field is a marshaling operation used to add a key-value pair to a logger's
	// context. Most fields are lazily marshaled, so it's inexpensive to add fields
	// to disabled debug-level log statements.
	Field = zapcore.Field

	// ObjectEncoder is a strongly-typed, encoding-agnostic interface for adding a
	// map- or struct-like object to the logging context. Like maps, ObjectEncoders
	// aren't safe for concurrent use (though typical use shouldn't require locks).
	ObjectEncoder = zapcore.ObjectEncoder

	// ObjectMarshaler allows user-defined types to efficiently add themselves to the
	// logging context, and to selectively omit information which shouldn't be
	// included in logs (e.g., passwords).
	//
	// Note: ObjectMarshaler is only used when Object is used or when
	// passed directly to Any. It is not used when reflection-based
	// encoding is used.
	ObjectMarshaler = zapcore.ObjectMarshaler
)

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

// Object constructs a field with the given key and ObjectMarshaler. It
// provides a flexible, but still type-safe and efficient, way to add map- or
// struct-like user-defined types to the logging context. The struct's
// MarshalLogObject method is called lazily.
func Object(key string, val ObjectMarshaler) Field {
	return Field{Key: key, Type: zapcore.ObjectMarshalerType, Interface: val}
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, value []byte) Field {
	return Field{Key: key, Type: zapcore.BinaryType, Interface: value}
}

// Bool constructs a field that carries a bool.
func Bool(key string, value bool) Field {
	var v int64
	if value {
		v = 1
	}
	return Field{Key: key, Type: zapcore.BoolType, Integer: v}
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: zapcore.ByteStringType, Interface: val}
}

// Namespace creates a named, isolated scope within the logger's context. All
// subsequent fields will be added to the new namespace.
//
// This helps prevent key collisions when injecting loggers into sub-components
// or third-party libraries.
func Namespace(key string) Field {
	return Field{Key: key, Type: zapcore.NamespaceType}
}

func Method(value string) Field {
	return Field{Key: `method`, Type: zapcore.StringType, String: value}
}

func Action(value string) Field {
	return Field{Key: `action`, Type: zapcore.StringType, String: value}
}

// Topic constructs a field that carries the name of the kafka topic.
func Topic(value string) Field {
	return Field{Key: `topic`, Type: zapcore.StringType, String: value}
}

// Partition constructs a field that carries the kafka partition number.
func Partition(value int) Field {
	return Field{Key: `partition`, Type: zapcore.Int64Type, Integer: int64(value)}
}

// Offset constructs a field that carries the value of the current kafka offset.
func Offset(value int64) Field {
	return Field{Key: `offset`, Type: zapcore.Int64Type, Integer: value}
}

func ProductID(value uint64) Field {
	return Field{Key: `product_id`, Type: zapcore.Uint64Type, Integer: int64(value)}
}

func Error(err error) Field {
	return Field{Key: `error`, Type: zapcore.StringType, String: err.Error()}
}

func Count(count int) Field {
	return Field{Key: `count`, Type: zapcore.Int64Type, Integer: int64(count)}
}

func Query(query string) Field {
	return Field{Key: `query`, Type: zapcore.StringType, String: query}
}

func File(fileName string) Field {
	return Field{Key: `file`, Type: zapcore.StringType, String: fileName}
}

const NoTraceId = `unknown`

// TraceId - extract trace ID from span
func TraceId(ctx context.Context) Field {
	span := trace.SpanFromContext(ctx)

	if span.SpanContext().TraceID().IsValid() {
		return Field{
			Key:    `traceId`,
			Type:   zapcore.StringType,
			String: span.SpanContext().TraceID().String(),
		}
	} else {
		return Field{
			Key:    `traceId`,
			Type:   zapcore.StringType,
			String: NoTraceId,
		}
	}
}
