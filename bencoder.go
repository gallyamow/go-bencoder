package bencoder

import (
	"bufio"
	"io"
)

func Encode(val any) []byte {
	en := encoder{}
	en.encodeUnknown(val)
	return en.bytes()
}

func Decode(rdr io.Reader) (any, error) {
	de := decoder{
		rdr: bufio.NewReader(rdr),
	}
	return de.decodeUnknown()
}
