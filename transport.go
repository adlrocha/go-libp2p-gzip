package gzip

import (
	"compress/gzip"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/compression"
	"go.uber.org/multierr"
)

// ID is the protocol ID for noise
const ID = "/compression/gzip"

var _ compression.CompressedTransport = &Transport{}

// Transport implements the interface sec.SecureTransport
// https://godoc.org/github.com/libp2p/go-libp2p-core/sec#SecureConn
type compConn struct {
	rlock sync.Mutex
	wlock sync.Mutex
	raw   net.Conn

	w *gzip.Writer
	r *gzip.Reader
}

// Transport defines a compression transport with a compression level.
type Transport struct {
	level int
}

// New Creates a new tranport with a specific compression level.
func New() *Transport {
	return &Transport{}
}

//NewConn upgrades a raw connection into a compressed connection.
func (t *Transport) NewConn(raw net.Conn, isServer bool) (compression.CompressedConn, error) {
	// if t.level > 0 && t.level <= 9 {
	// 	w, err := gzip.NewWriterLevel(raw, t.level)
	// 	return &compConn{
	// 		raw: raw,
	// 		w:   w,
	// 	}, err
	// }
	return &compConn{
		raw: raw,
		w:   gzip.NewWriter(raw),
	}, nil
}

// Write compression wrapper
func (c *compConn) Write(b []byte) (int, error) {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	n, err := c.w.Write(b)
	return n, multierr.Combine(err, c.w.Flush())
}

// Read compression wrapper
func (c *compConn) Read(b []byte) (int, error) {
	c.rlock.Lock()
	defer c.rlock.Unlock()
	if c.r == nil {
		// This _needs_ to be lazy as it reads a header.
		var err error
		c.r, err = gzip.NewReader(c.raw)
		if err != nil {
			return 0, err
		}
	}
	n, err := c.r.Read(b)
	if err != nil {
		c.r.Close()
	}
	return n, err
}

func (c *compConn) Close() error {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	return multierr.Combine(c.w.Close(), c.raw.Close())
}

func (c *compConn) LocalAddr() net.Addr {
	return c.raw.LocalAddr()
}

func (c *compConn) RemoteAddr() net.Addr {
	return c.raw.RemoteAddr()
}

func (c *compConn) SetDeadline(t time.Time) error {
	return c.raw.SetDeadline(t)
}

func (c *compConn) SetReadDeadline(t time.Time) error {
	return c.raw.SetReadDeadline(t)
}

func (c *compConn) SetWriteDeadline(t time.Time) error {
	return c.raw.SetWriteDeadline(t)
}
