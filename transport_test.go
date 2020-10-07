package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"math/rand"
	"net"
	"testing"

	"github.com/libp2p/go-libp2p-core/compression"
)

func newTestTransport(t *testing.T, typ, bits int) *Transport {
	return New(6)
}

// Create a new pair of connected TCP sockets.
func newConnPair(t *testing.T) (net.Conn, net.Conn) {
	lstnr, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
		return nil, nil
	}

	var clientErr error
	var client net.Conn
	addr := lstnr.Addr()
	done := make(chan struct{})

	go func() {
		defer close(done)
		client, clientErr = net.Dial(addr.Network(), addr.String())
	}()

	server, err := lstnr.Accept()
	<-done

	lstnr.Close()

	if err != nil {
		t.Fatalf("Failed to accept: %v", err)
	}

	if clientErr != nil {
		t.Fatalf("Failed to connect: %v", clientErr)
	}

	return client, server
}

//Test a compressed exchange
func TestExchange(t *testing.T) {
	init, resp := newConnPair(t)
	initTransport := New(gzip.DefaultCompression)
	respTransport := New(6)

	var initConn, respConn compression.CompressedConn

	initConn, _ = initTransport.NewConn(init, true)
	respConn, _ = respTransport.NewConn(resp, false)

	defer initConn.Close()
	defer respConn.Close()

	size := 100000

	before := makeLargePlaintext(size)
	_, err := initConn.Write(before)
	if err != nil {
		t.Fatal(err)
	}

	after := make([]byte, len(before))
	afterLen, err := io.ReadFull(respConn, after)
	if err != nil {
		t.Fatal(err)
	}

	if len(before) != afterLen {
		t.Errorf("expected to read same amount of data as written. written=%d read=%d", len(before), afterLen)
	}
	if !bytes.Equal(before, after) {
		t.Error("Message mismatch.")
	}
}

// Test that the connection is compressed by checking a mismatch
// between a compressed and uncompressed connection.
func TestMismatch(t *testing.T) {
	init, respConn := newConnPair(t)
	initTransport := New(gzip.DefaultCompression)

	var initConn compression.CompressedConn

	initConn, _ = initTransport.NewConn(init, true)

	defer initConn.Close()
	defer respConn.Close()

	size := 100000

	before := makeLargePlaintext(size)
	_, err := initConn.Write(before)
	if err != nil {
		t.Fatal(err)
	}

	after := make([]byte, len(before))
	_, err = io.ReadFull(respConn, after)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(before, after) {
		t.Error("The received message should have been compressed because I am using a raw connection.")
	}
}

// Benchmark the use of compression with large plaintext.
func TestCompressedBenchmark(t *testing.T) {
	init, resp := newConnPair(t)
	initTransport := New(gzip.DefaultCompression)
	respTransport := New(6)

	var initConn, respConn compression.CompressedConn

	initConn, _ = initTransport.NewConn(init, true)
	respConn, _ = respTransport.NewConn(resp, false)

	defer initConn.Close()
	defer respConn.Close()

	size := 100000
	before := makeLargePlaintext(size)
	for i := 0; i < 1000; i++ {
		_, err := initConn.Write(before)
		if err != nil {
			t.Fatal(err)
		}

		after := make([]byte, len(before))
		afterLen, err := io.ReadFull(respConn, after)
		if err != nil {
			t.Fatal(err)
		}
		if len(before) != afterLen {
			t.Errorf("expected to read same amount of data as written. written=%d read=%d", len(before), afterLen)
		}
		if !bytes.Equal(before, after) {
			t.Error("Message mismatch.")
		}
	}

}

func makeLargePlaintext(size int) []byte {
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}
