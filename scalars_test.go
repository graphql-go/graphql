package graphql

import (
	"math"
	"testing"
)

func TestCoerceInt(t *testing.T) {
	tests := []struct {
		in   interface{}
		want interface{}
	}{
		{
			in:   false,
			want: 0,
		},
		{
			in:   true,
			want: 1,
		},
		{
			in:   boolPtr(false),
			want: 0,
		},
		{
			in:   boolPtr(true),
			want: 1,
		},
		{
			in:   (*bool)(nil),
			want: nil,
		},
		{
			in:   int(math.MinInt32) - 1,
			want: nil,
		},
		{
			in:   int(math.MaxInt32) + 1,
			want: nil,
		},
		{
			in:   uint(math.MaxInt32) + 1,
			want: nil,
		},
		{
			in:   uint32(math.MaxInt32) + 1,
			want: nil,
		},
		{
			in:   int64(math.MinInt32) - 1,
			want: nil,
		},
		{
			in:   int64(math.MaxInt32) + 1,
			want: nil,
		},
		{
			in:   uint64(math.MaxInt32) + 1,
			want: nil,
		},
		{
			// need to subtract more than one because of float32 precision
			in:   float32(math.MinInt32) - 1000,
			want: nil,
		},
		{
			// need to add more than one because of float32 precision
			in:   float32(math.MaxInt32) + 1000,
			want: nil,
		},
		{
			in:   float64(math.MinInt32) - 1,
			want: nil,
		},
		{
			in:   float64(math.MaxInt32) + 1,
			want: nil,
		},
		{
			in:   int(math.MinInt32),
			want: int(math.MinInt32),
		},
		{
			in:   int(math.MaxInt32),
			want: int(math.MaxInt32),
		},
		{
			in:   intPtr(12),
			want: 12,
		},
		{
			in:   (*int)(nil),
			want: nil,
		},
		{
			in:   int8(13),
			want: int(13),
		},
		{
			in:   int8Ptr(14),
			want: int(14),
		},
		{
			in:   (*int8)(nil),
			want: nil,
		},
		{
			in:   int16(15),
			want: int(15),
		},
		{
			in:   int16Ptr(16),
			want: int(16),
		},
		{
			in:   (*int16)(nil),
			want: nil,
		},
		{
			in:   int32(17),
			want: int(17),
		},
		{
			in:   int32Ptr(18),
			want: int(18),
		},
		{
			in:   (*int32)(nil),
			want: nil,
		},
		{
			in:   int64(19),
			want: int(19),
		},
		{
			in:   int64Ptr(20),
			want: int(20),
		},
		{
			in:   (*int64)(nil),
			want: nil,
		},
		{
			in:   uint8(21),
			want: int(21),
		},
		{
			in:   uint8Ptr(22),
			want: int(22),
		},
		{
			in:   (*uint8)(nil),
			want: nil,
		},
		{
			in:   uint16(23),
			want: int(23),
		},
		{
			in:   uint16Ptr(24),
			want: int(24),
		},
		{
			in:   (*uint16)(nil),
			want: nil,
		},
		{
			in:   uint32(25),
			want: int(25),
		},
		{
			in:   uint32Ptr(26),
			want: int(26),
		},
		{
			in:   (*uint32)(nil),
			want: nil,
		},
		{
			in:   uint64(27),
			want: int(27),
		},
		{
			in:   uint64Ptr(28),
			want: int(28),
		},
		{
			in:   (*uint64)(nil),
			want: nil,
		},
		{
			in:   uintPtr(29),
			want: int(29),
		},
		{
			in:   (*uint)(nil),
			want: nil,
		},
		{
			in:   float32(30.1),
			want: int(30),
		},
		{
			in:   float32Ptr(31.2),
			want: int(31),
		},
		{
			in:   (*float32)(nil),
			want: nil,
		},
		{
			in:   float64(32),
			want: int(32),
		},
		{
			in:   float64Ptr(33.1),
			want: int(33),
		},
		{
			in:   (*float64)(nil),
			want: nil,
		},
		{
			in:   "34",
			want: int(34),
		},
		{
			in:   stringPtr("35"),
			want: int(35),
		},
		{
			in:   (*string)(nil),
			want: nil,
		},
		{
			in:   "I'm not a number",
			want: nil,
		},
		{
			in:   make(map[string]interface{}),
			want: nil,
		},
	}

	for i, tt := range tests {
		if got, want := coerceInt(tt.in), tt.want; got != want {
			t.Errorf("%d: in=%v, got=%v, want=%v", i, tt.in, got, want)
		}
	}
}

func TestCoerceFloat(t *testing.T) {
	tests := []struct {
		in   interface{}
		want interface{}
	}{
		{
			in:   false,
			want: 0.0,
		},
		{
			in:   true,
			want: 1.0,
		},
		{
			in:   boolPtr(false),
			want: 0.0,
		},
		{
			in:   boolPtr(true),
			want: 1.0,
		},
		{
			in:   (*bool)(nil),
			want: nil,
		},
		{
			in:   int(math.MinInt32),
			want: float64(math.MinInt32),
		},
		{
			in:   int(math.MaxInt32),
			want: float64(math.MaxInt32),
		},
		{
			in:   intPtr(12),
			want: float64(12),
		},
		{
			in:   (*int)(nil),
			want: nil,
		},
		{
			in:   int8(13),
			want: float64(13),
		},
		{
			in:   int8Ptr(14),
			want: float64(14),
		},
		{
			in:   (*int8)(nil),
			want: nil,
		},
		{
			in:   int16(15),
			want: float64(15),
		},
		{
			in:   int16Ptr(16),
			want: float64(16),
		},
		{
			in:   (*int16)(nil),
			want: nil,
		},
		{
			in:   int32(17),
			want: float64(17),
		},
		{
			in:   int32Ptr(18),
			want: float64(18),
		},
		{
			in:   (*int32)(nil),
			want: nil,
		},
		{
			in:   int64(19),
			want: float64(19),
		},
		{
			in:   int64Ptr(20),
			want: float64(20),
		},
		{
			in:   (*int64)(nil),
			want: nil,
		},
		{
			in:   uint8(21),
			want: float64(21),
		},
		{
			in:   uint8Ptr(22),
			want: float64(22),
		},
		{
			in:   (*uint8)(nil),
			want: nil,
		},
		{
			in:   uint16(23),
			want: float64(23),
		},
		{
			in:   uint16Ptr(24),
			want: float64(24),
		},
		{
			in:   (*uint16)(nil),
			want: nil,
		},
		{
			in:   uint32(25),
			want: float64(25),
		},
		{
			in:   uint32Ptr(26),
			want: float64(26),
		},
		{
			in:   (*uint32)(nil),
			want: nil,
		},
		{
			in:   uint64(27),
			want: float64(27),
		},
		{
			in:   uint64Ptr(28),
			want: float64(28),
		},
		{
			in:   (*uint64)(nil),
			want: nil,
		},
		{
			in:   uintPtr(29),
			want: float64(29),
		},
		{
			in:   (*uint)(nil),
			want: nil,
		},
		{
			in:   float32(30),
			want: float32(30),
		},
		{
			in:   float32Ptr(31),
			want: float32(31),
		},
		{
			in:   (*float32)(nil),
			want: nil,
		},
		{
			in:   float64(32),
			want: float64(32),
		},
		{
			in:   float64Ptr(33.2),
			want: float64(33.2),
		},
		{
			in:   (*float64)(nil),
			want: nil,
		},
		{
			in:   "34",
			want: float64(34),
		},
		{
			in:   stringPtr("35.2"),
			want: float64(35.2),
		},
		{
			in:   (*string)(nil),
			want: nil,
		},
		{
			in:   "I'm not a number",
			want: nil,
		},
		{
			in:   make(map[string]interface{}),
			want: nil,
		},
	}

	for i, tt := range tests {
		if got, want := coerceFloat(tt.in), tt.want; got != want {
			t.Errorf("%d: in=%v, got=%v, want=%v", i, tt.in, got, want)
		}
	}
}

func TestCoerceBool(t *testing.T) {
	tests := []struct {
		in   interface{}
		want interface{}
	}{
		{
			in:   false,
			want: false,
		},
		{
			in:   true,
			want: true,
		},
		{
			in:   boolPtr(false),
			want: false,
		},
		{
			in:   boolPtr(true),
			want: true,
		},
		{
			in:   (*bool)(nil),
			want: nil,
		},
		{
			in:   int(math.MinInt32),
			want: true,
		},
		{
			in:   int(math.MaxInt32),
			want: true,
		},
		{
			in:   int(0),
			want: false,
		},
		{
			in:   intPtr(12),
			want: true,
		},
		{
			in:   intPtr(0),
			want: false,
		},
		{
			in:   (*int)(nil),
			want: nil,
		},
		{
			in:   int8(13),
			want: true,
		},
		{
			in:   int8(0),
			want: false,
		},
		{
			in:   int8Ptr(14),
			want: true,
		},
		{
			in:   int8Ptr(0),
			want: false,
		},
		{
			in:   (*int8)(nil),
			want: nil,
		},
		{
			in:   int16(15),
			want: true,
		},
		{
			in:   int16(0),
			want: false,
		},
		{
			in:   int16Ptr(16),
			want: true,
		},
		{
			in:   int16Ptr(0),
			want: false,
		},
		{
			in:   (*int16)(nil),
			want: nil,
		},
		{
			in:   int32(17),
			want: true,
		},
		{
			in:   int32(0),
			want: false,
		},
		{
			in:   int32Ptr(18),
			want: true,
		},
		{
			in:   int32Ptr(0),
			want: false,
		},
		{
			in:   (*int32)(nil),
			want: nil,
		},
		{
			in:   int64(19),
			want: true,
		},
		{
			in:   int64(0),
			want: false,
		},
		{
			in:   int64Ptr(20),
			want: true,
		},
		{
			in:   int64Ptr(0),
			want: false,
		},
		{
			in:   (*int64)(nil),
			want: nil,
		},
		{
			in:   uint8(21),
			want: true,
		},
		{
			in:   uint8(0),
			want: false,
		},
		{
			in:   uint8Ptr(22),
			want: true,
		},
		{
			in:   uint8Ptr(0),
			want: false,
		},
		{
			in:   (*uint8)(nil),
			want: nil,
		},
		{
			in:   uint16(23),
			want: true,
		},
		{
			in:   uint16(0),
			want: false,
		},
		{
			in:   uint16Ptr(24),
			want: true,
		},
		{
			in:   uint16Ptr(0),
			want: false,
		},
		{
			in:   (*uint16)(nil),
			want: nil,
		},
		{
			in:   uint32(25),
			want: true,
		},
		{
			in:   uint32(0),
			want: false,
		},
		{
			in:   uint32Ptr(26),
			want: true,
		},
		{
			in:   uint32Ptr(0),
			want: false,
		},
		{
			in:   (*uint32)(nil),
			want: nil,
		},
		{
			in:   uint64(27),
			want: true,
		},
		{
			in:   uint64(0),
			want: false,
		},
		{
			in:   uint64Ptr(28),
			want: true,
		},
		{
			in:   uint64Ptr(0),
			want: false,
		},
		{
			in:   (*uint64)(nil),
			want: nil,
		},
		{
			in:   uintPtr(29),
			want: true,
		},
		{
			in:   uintPtr(0),
			want: false,
		},
		{
			in:   float32(30),
			want: true,
		},
		{
			in:   float32(0),
			want: false,
		},
		{
			in:   float32Ptr(31),
			want: true,
		},
		{
			in:   float32Ptr(0),
			want: false,
		},
		{
			in:   (*float32)(nil),
			want: nil,
		},
		{
			in:   float64(32),
			want: true,
		},
		{
			in:   float64(0),
			want: false,
		},
		{
			in:   float64Ptr(33.2),
			want: true,
		},
		{
			in:   float64Ptr(0),
			want: false,
		},
		{
			in:   (*float64)(nil),
			want: nil,
		},
		{
			in:   "34",
			want: true,
		},
		{
			in:   "false",
			want: false,
		},
		{
			in:   stringPtr("true"),
			want: true,
		},
		{
			in:   stringPtr("false"),
			want: false,
		},
		{
			in:   (*string)(nil),
			want: nil,
		},
		{
			in:   "I'm some random string",
			want: true,
		},
		{
			in:   "",
			want: false,
		},
		{
			in:   int8(0),
			want: false,
		},
		{
			in:   make(map[string]interface{}),
			want: false,
		},
	}

	for i, tt := range tests {
		if got, want := coerceBool(tt.in), tt.want; got != want {
			t.Errorf("%d: in=%v, got=%v, want=%v", i, tt.in, got, want)
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(n int) *int {
	return &n
}

func int8Ptr(n int8) *int8 {
	return &n
}

func int16Ptr(n int16) *int16 {
	return &n
}

func int32Ptr(n int32) *int32 {
	return &n
}

func int64Ptr(n int64) *int64 {
	return &n
}

func uintPtr(n uint) *uint {
	return &n
}

func uint8Ptr(n uint8) *uint8 {
	return &n
}

func uint16Ptr(n uint16) *uint16 {
	return &n
}

func uint32Ptr(n uint32) *uint32 {
	return &n
}

func uint64Ptr(n uint64) *uint64 {
	return &n
}

func float32Ptr(n float32) *float32 {
	return &n
}

func float64Ptr(n float64) *float64 {
	return &n
}

func stringPtr(s string) *string {
	return &s
}
