package bench

import (
	"bencoder"
	"bytes"
	jackpal "github.com/jackpal/bencode-go"
	"strings"
	"testing"
)

func BenchmarkBencoderEncode(b *testing.B) {
	mp := map[string]interface{}{
		"key": "value",
		"num": 123,
	}
	array := []interface{}{
		123,
		"string",
		map[string]interface{}{
			"key": "value",
		},
	}

	for b.Loop() {
		bencoder.Encode(mp)
		bencoder.Encode(array)
	}
}

func BenchmarkJackpalEncode(b *testing.B) {
	dict := map[string]interface{}{
		"key": "value",
		"num": 123,
	}
	array := []interface{}{
		123,
		"string",
		map[string]interface{}{
			"key": "value",
		},
	}

	for b.Loop() {
		var buf bytes.Buffer
		if err := jackpal.Marshal(&buf, dict); err != nil {
			b.Fatal(err)
		}
		if err := jackpal.Marshal(&buf, array); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBencoderDecode(b *testing.B) {
	dict := "d3:key5:value3:numi123ee"
	array := "li123e6:stringd3:key5:valueee"

	for b.Loop() {
		if _, err := bencoder.Decode(strings.NewReader(dict)); err != nil {
			b.Fatal(err)
		}
		if _, err := bencoder.Decode(strings.NewReader(array)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJackpalDecode(b *testing.B) {
	dict := "d3:key5:value3:numi123ee"
	array := "li123e6:stringd3:key5:valueee"

	for b.Loop() {
		if _, err := jackpal.Decode(strings.NewReader(dict)); err != nil {
			b.Fatal(err)
		}
		if _, err := jackpal.Decode(strings.NewReader(array)); err != nil {
			b.Fatal(err)
		}
	}
}
