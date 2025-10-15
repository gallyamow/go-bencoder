package bencoder

func Encode(val any) []byte {
	en := encoder{}
	en.encodeUnknown(val)
	return en.bytes()
}

func Decode(s string) (map[string]any, error) {
	return nil, nil
}
