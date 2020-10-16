module github.com/libp2p/go-libp2p-gzip

go 1.14

// Compression with go-libp2p-core 0.7
// replace github.com/libp2p/go-libp2p-core => github.com/adlrocha/go-libp2p-core bf5d45ca7e53c6c20d8f319f0f23c2fef93bf7f6
// Compression with go-libp2p-core 0.6
replace github.com/libp2p/go-libp2p-core => github.com/adlrocha/go-libp2p-core v0.6.2-0.20201008055055-b309947fc237

// replace github.com/libp2p/go-libp2p-core => github.com/adlrocha/go-libp2p-core v0.6.2-0.20201007141150-bf7ae45bb37e

require (
	github.com/klauspost/compress v1.11.1
	github.com/libp2p/go-libp2p-core v0.0.0-00010101000000-000000000000
	go.uber.org/multierr v1.6.0
)
