## go-bencoder

Simple `bencode` encoder/decoder implementation. No reflection usage.

* https://ru.wikipedia.org/wiki/Bencode

### Usage

```go
package main

import (
	"fmt"
	"strings"

	bencoder "github.com/gallyamow/go-bencoder"
)

func main() {
	// decoding
	dict, err := bencoder.Decode(strings.NewReader("d2:k1d1:ai1e1:bi2e1:ci3e1:di4e1:ei5ee2:k2d1:a1:z1:b1:x1:c1:y1:d1:q1:e1:pee"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded dict: %+v\n", dict)

	// encoding
	s := bencoder.Encode(map[string]any{
		"k1": map[string]any{"a": 1, "b": int8(2), "c": int16(3), "d": int32(4), "e": int64(5)},
		"k2": map[string]any{"e": "p", "d": "q", "c": "y", "b": "x", "a": "z"},
	})
	fmt.Printf("Encoded dict: %+v\n", string(s))
}
```
