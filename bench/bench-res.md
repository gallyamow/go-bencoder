```
goos: darwin
goarch: arm64
pkg: bencoder/bench
cpu: Apple M4
BenchmarkBencoderEncode-10    	 4320560	       257.3 ns/op	     246 B/op	       9 allocs/op
BenchmarkJackpalEncode-10     	 1482273	       809.6 ns/op	     592 B/op	      21 allocs/op
BenchmarkBencoderDecode-10    	 1000000	       1067 ns/op	    9440 B/op	      37 allocs/op
BenchmarkJackpalDecode-10     	 2173898	       563.0 ns/op	    1021 B/op	      28 allocs/op
PASS
ok  	bencoder/bench	4.913s
```