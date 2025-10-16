package bencoder

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

var stringTests = []struct {
	name    string
	raw     string
	encoded string
}{
	{"ascii", "spam", "4:spam"},
	{"ascii with space", "with space", "10:with space"},
	{"russian", "на русском", "19:на русском"},
	{"chinese", "拉面", "6:拉面"},
}

var intTests = []struct {
	name    string
	encoded string
	raw     int64
}{
	{"positive small value", "i1234e", 1234},
	{"zero", "i0e", 0},
	{"negative", "i-1234e", -1234},
	{"negative border", "i-9223372036854775808e", math.MinInt64},
	{"positive border", "i9223372036854775807e", math.MaxInt64},
}

var uintTests = []struct {
	name    string
	encoded string
	raw     uint64
}{
	{"pre border value (MaxInt64+1)", "i9223372036854775808e", uint64(9223372036854775808)},
	{"positive border", "i18446744073709551615e", math.MaxUint64},
}

var arrayTests = []struct {
	name    string
	raw     []any
	encoded string
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

var dictTests = []struct {
	name    string
	raw     map[string]any
	encoded string
}{
	{"map of strings", map[string]any{"k1": "val1", "k2": "val2"}, "d2:k14:val12:k24:val2e"},
	{"map of int", map[string]any{"k1": int8(1), "k2": uint8(2)}, "d2:k1i1e2:k2i2ee"},
	{"map of arrays of int", map[string]any{"k1": []int{1, 2, 3}, "k1234": []uint{1, 2, 3}}, "d2:k1li1ei2ei3ee5:k1234li1ei2ei3eee"},
	{"map of map[string]any", map[string]any{
		"k1": map[string]any{"a": 1, "b": int8(2), "c": int16(3), "d": int32(4), "e": int64(5)},
		"k2": map[string]any{"e": "p", "d": "q", "c": "y", "b": "x", "a": "z"},
	}, "d2:k1d1:ai1e1:bi2e1:ci3e1:di4e1:ei5ee2:k2d1:a1:z1:b1:x1:c1:y1:d1:q1:e1:pee"},
}

func TestEncode(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		for _, tt := range stringTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.raw)

				if string(got) != tt.encoded {
					t.Errorf("Encode() = %q, want %q", got, tt.encoded)
				}
			})
		}
	})

	t.Run("Integer", func(t *testing.T) {
		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.raw)

				if string(got) != tt.encoded {
					t.Errorf("Encode() = %q, want %q", got, tt.encoded)
				}
			})
		}

		for _, tt := range uintTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.raw)

				if string(got) != tt.encoded {
					t.Errorf("Encode() = %q, want %q", got, tt.encoded)
				}
			})
		}
	})

	t.Run("Array", func(t *testing.T) {
		for _, tt := range arrayTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.raw)

				if string(got) != tt.encoded {
					t.Errorf("Encode() = %q, want %q", got, tt.encoded)
				}
			})
		}
	})

	t.Run("Dict", func(t *testing.T) {
		for _, tt := range dictTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Encode(tt.raw)

				if string(got) != tt.encoded {
					t.Errorf("Encode() = %q, want %q", got, tt.encoded)
				}
			})
		}
	})
}

func TestDecode(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		for _, tt := range stringTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Decode(strings.NewReader(tt.encoded))

				if err != nil {
					t.Errorf("Decode() error = %v", err)
				}

				if got.(string) != tt.raw {
					t.Errorf("Encode() = %q, want %q", got, tt.raw)
				}
			})
		}
	})

	t.Run("Integer", func(t *testing.T) {
		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Decode(strings.NewReader(tt.encoded))

				if err != nil {
					t.Fatalf("Decode() error = %v", err)
				}

				if got.(int64) != tt.raw {
					t.Errorf("Decode() = %d, want %d", got, tt.raw)
				}
			})
		}

		for _, tt := range uintTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Decode(strings.NewReader(tt.encoded))

				if err != nil {
					t.Fatalf("Decode() error = %v", err)
				}

				if got.(uint64) != tt.raw {
					t.Errorf("Decode() = %d, want %d", got, tt.raw)
				}
			})
		}
	})

	t.Run("Array", func(t *testing.T) {
		for _, tt := range arrayTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Decode(strings.NewReader(tt.encoded))

				if err != nil {
					t.Fatalf("Decode() error = %v", err)
				}

				// cant use DeepEqual due different int types
				// For more robust type-tolerant equality, you could normalize recursively or use cmpopts.IgnoreTypes from google/go-cmp
				for i, v := range got.([]any) {
					// can't check equality without normalization, so use this simple way
					if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", tt.raw[i]) {
						t.Errorf("Decode() = got[%T] %#v, want[%T] %#v", v, v, tt.raw[i], tt.raw[i])
					}
				}
			})
		}
	})

	t.Run("Dict", func(t *testing.T) {
		for _, tt := range dictTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Decode(strings.NewReader(tt.encoded))

				if err != nil {
					t.Fatalf("Decode() error = %v", err)
				}

				// cant use DeepEqual due different int types
				for k, v := range got.(map[string]any) {
					// can't check equality without normalization, so use this simple way
					if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", tt.raw[k]) {
						t.Errorf("Decode() = got[%T] %#v, want[%T] %#v", v, v, tt.raw[k], tt.raw[k])
					}
				}
			})
		}
	})
}
