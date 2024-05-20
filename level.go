package zwrap

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Level interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

func castToZlogLevel(level any) zerolog.Level {
	switch casted := level.(type) {
	case int:
		return toZlogLevel[int](casted)
	case int8:
		return toZlogLevel[int8](casted)
	case int16:
		return toZlogLevel[int16](casted)
	case int32:
		return toZlogLevel[int32](casted)
	case int64:
		return toZlogLevel[int64](casted)
	case uint:
		return toZlogLevel[uint](casted)
	case uint8:
		return toZlogLevel[uint8](casted)
	case uint16:
		return toZlogLevel[uint16](casted)
	case uint32:
		return toZlogLevel[uint32](casted)
	case uint64:
		return toZlogLevel[uint64](casted)
	case string:
		if parsed, err := zerolog.ParseLevel(casted); err == nil {
			return parsed
		} else {
			panic(fmt.Sprintf("invalid log level string %v: %v", level, err))
		}
	case zerolog.Level:
		return casted
	default:
		panic(fmt.Sprintf("invalid log level type (%T): %v", level, level))
	}
}

func toZlogLevel[T Level](level T) zerolog.Level {
	switch casted := any(level).(type) {
	case uint32: // compat
		switch {
		case casted == 0:
			return zerolog.PanicLevel
		case casted == 1:
			return zerolog.FatalLevel
		case casted == 2:
			return zerolog.ErrorLevel
		case casted == 3:
			return zerolog.WarnLevel
		case casted == 4:
			return zerolog.InfoLevel
		case casted == 5:
			return zerolog.DebugLevel
		case casted == 6:
			return zerolog.TraceLevel
		default:
			if casted < 0 {
				return zerolog.TraceLevel
			}
			if casted > 5 {
				return zerolog.PanicLevel
			}
		}
	case int16, int32, int64, int, uint, uint8, uint16, uint64, int8:
		if level < 0 {
			return zerolog.TraceLevel
		}
		if level > 5 {
			return zerolog.PanicLevel
		}
		return zerolog.Level(int8(level))
	}
	return zerolog.TraceLevel
}
