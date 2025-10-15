package bencoder

import (
	"math"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  string
		}{
			{"ascii", "spam", "4:spam"},
			{"ascii with space", "with space", "10:with space"},
			{"russian", "на русском", "19:на русском"},
			{"chinese", "拉面", "6:拉面"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.input)

				if string(got) != tt.want {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("Integer", func(t *testing.T) {
		testsInt := []struct {
			name  string
			input int64
			want  string
		}{
			{"positive", 1234, "i1234e"},
			{"zero", 0, "i0e"},
			{"negative", -1234, "i-1234e"},
			{"longlong", math.MaxInt64, "i9223372036854775807e"},
		}

		for _, tt := range testsInt {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.input)

				if string(got) != tt.want {
					t.Errorf("Encode() = %q, want %s", got, tt.want)
				}
			})
		}

		testsUint := []struct {
			name  string
			input uint64
			want  string
		}{
			{"positive", 1234, "i1234e"},
			{"zero", 0, "i0e"},
			{"longlong", math.MaxUint64, "i18446744073709551615e"},
		}

		for _, tt := range testsUint {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.input)

				if string(got) != tt.want {
					t.Errorf("Encode() = %q, want %s", got, tt.want)
				}
			})
		}
	})

	t.Run("Array", func(t *testing.T) {
		tests := []struct {
			name  string
			input []any
			want  string
		}{
			{"array of strings", []any{"spam", "spame", "spamer"}, "l4:spam5:spame6:spamere"},
			{"array of int arrays", []any{1, 2, 3}, "li1ei2ei3ee"},
			{"array of any int arrays", []any{[]int{1, 2, 3}, []int8{4, 5, 6}, []int16{7, 8, 9},
				[]int32{1, 2, 3}, []int64{4, 5, 6}}, "lli1ei2ei3eeli4ei5ei6eeli7ei8ei9eeli1ei2ei3eeli4ei5ei6eee"},
			{"array of maps[string]any", []any{
				map[string]any{"a": 1, "b": int8(2), "c": int16(3), "d": int32(4), "e": int64(5)}, // sorted keys
				map[string]any{"e": "p", "d": "q", "c": "y", "b": "x", "a": "z"},                  // unsorted keys
			}, "ld1:ai1e1:bi2e1:ci3e1:di4e1:ei5eed1:a1:z1:b1:x1:c1:y1:d1:q1:e1:pee"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.input)

				if string(got) != tt.want {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("Dict", func(t *testing.T) {
		tests := []struct {
			name  string
			input map[string]any
			want  string
		}{
			{"map of strings", map[string]any{"k1": "val1", "k2": "val2"}, "d2:k14:val12:k24:val2e"},
			{"map of int", map[string]any{"k1": int8(1), "k2": uint8(2)}, "d2:k1i1e2:k2i2ee"},
			{"map of arrays of int", map[string]any{"k1": []int{1, 2, 3}, "k1234": []uint{1, 2, 3}}, "d2:k1li1ei2ei3ee5:k1234li1ei2ei3eee"},
			{"map of map[string]any", map[string]any{
				"k1": map[string]any{"a": 1, "b": int8(2), "c": int16(3), "d": int32(4), "e": int64(5)},
				"k2": map[string]any{"e": "p", "d": "q", "c": "y", "b": "x", "a": "z"},
			}, "d2:k1d1:ai1e1:bi2e1:ci3e1:di4e1:ei5ee2:k2d1:a1:z1:b1:x1:c1:y1:d1:q1:e1:pee"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.input)

				if string(got) != tt.want {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			})
		}
	})
}
