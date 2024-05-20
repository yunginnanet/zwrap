package zwrap

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestCastToZlogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level any
		want  zerolog.Level
	}{
		{"int", 3, zerolog.ErrorLevel},
		{"int8", int8(4), zerolog.FatalLevel},
		{"int16", int16(2), zerolog.WarnLevel},
		{"int32", int32(1), zerolog.InfoLevel},
		{"int64", int64(5), zerolog.PanicLevel},
		{"uint", uint(6), zerolog.PanicLevel},
		{"uint8", uint8(0), zerolog.DebugLevel},
		{"uint16", uint16(4), zerolog.FatalLevel},
		{"uint32", uint32(3), zerolog.WarnLevel},
		{"uint64", uint64(2), zerolog.WarnLevel},
		{"string", "info", zerolog.InfoLevel},
		{"zerolog.Level", zerolog.DebugLevel, zerolog.DebugLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := castToZlogLevel(tt.level); got != tt.want {
				t.Errorf("castToZlogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastToZlogLevel_Panic(t *testing.T) {
	tests := []struct {
		name  string
		level any
	}{
		{"invalid type", 3.14},
		{"invalid string", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("castToZlogLevel() did not panic")
				}
			}()
			castToZlogLevel(tt.level)
		})
	}
}

func TestToZlogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level any
		want  zerolog.Level
	}{
		{"int8", int8(3), zerolog.ErrorLevel},
		{"int16", int16(2), zerolog.WarnLevel},
		{"int32", int32(1), zerolog.InfoLevel},
		{"int64", int64(5), zerolog.PanicLevel},
		{"uint", uint(0), zerolog.DebugLevel},
		{"uint8", uint8(6), zerolog.PanicLevel},
		{"uint16", uint16(4), zerolog.FatalLevel},
		{"uint32", uint32(3), zerolog.WarnLevel},
		{"uint64", uint64(2), zerolog.WarnLevel},
		{"-1 trace level", -1, zerolog.TraceLevel},
		{"> limit level", uint(7), zerolog.PanicLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch casted := tt.level.(type) {
			case int8:
				if got := toZlogLevel[int8](casted); got != tt.want {
					t.Errorf("toZlogLevel(%v) = %v, want %v", tt.level, got, tt.want)
				}
			case int16:
				if got := toZlogLevel[int16](casted); got != tt.want {
					t.Errorf("toZlogLevel(%v) = %v, want %v", tt.level, got, tt.want)
				}
			case int32:
				if got := toZlogLevel[int32](casted); got != tt.want {
					t.Errorf("toZlogLevel(%v) = %v, want %v", tt.level, got, tt.want)
				}
			case int64:
				if got := toZlogLevel[int64](casted); got != tt.want {
					t.Errorf("toZlogLevel(%v) = %v, want %v", tt.level, got, tt.want)
				}
			case uint:
				if got := toZlogLevel[uint](casted); got != tt.want {
					t.Errorf("toZlogLevel(%v) = %v, want %v", tt.level, got, tt.want)
				}
			case uint8:
				if got := toZlogLevel[uint8](casted); got != tt.want {
					t.Errorf("toZlogLevel(%v) = %v, want %v", tt.level, got, tt.want)
				}
			}
		})
	}
}
