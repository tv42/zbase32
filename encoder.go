package zbase32

import "io"

type encoder struct {
	io.WriteCloser
	w    io.Writer
	buf  [5]byte    // buffered data waiting to be encoded
	nbuf int        // number of bytes in buf
	out  [1024]byte // output buffer
	err  error
}

func (e *encoder) Write(p []byte) (n int, err error) {
	if e.err != nil {
		return 0, e.err
	}

	// Leading fringe.
	if e.nbuf > 0 {
		var i int
		for i = 0; i < len(p) && e.nbuf < 5; i++ {
			e.buf[e.nbuf] = p[i]
			e.nbuf++
		}
		n += i
		p = p[i:]
		if e.nbuf < 5 {
			return
		}
		m := encode(e.out[0:], e.buf[0:], -1)
		if _, e.err = e.w.Write(e.out[0:m]); e.err != nil {
			return n, e.err
		}
		e.nbuf = 0
	}

	// Large interior chunks.
	for len(p) >= 5 {
		nn := len(e.out) / 8 * 5
		if nn > len(p) {
			nn = len(p)
			nn -= nn % 5
		}
		m := encode(e.out[0:], p[0:nn], -1)
		if _, e.err = e.w.Write(e.out[0:m]); e.err != nil {
			return n, e.err
		}
		n += nn
		p = p[nn:]
	}

	// Trailing fringe.
	for i := 0; i < len(p); i++ {
		e.buf[i] = p[i]
	}
	e.nbuf = len(p)
	n += len(p)
	return
}

// Close flushes any pending output from the encoder. It is an error to call
// Write after calling Close.
func (e *encoder) Close() error {
	// If there's anything left in the buffer, flush it out
	if e.err == nil && e.nbuf > 0 {
		m := encode(e.out[0:], e.buf[0:e.nbuf], -1)
		_, e.err = e.w.Write(e.out[0:m])
		e.nbuf = 0
	}
	return e.err
}

// NewEncoder returns a new z-base-32 stream encoder. Data written to the
// returned writer will be encoded and then written to w. z-Base-32
// encodings operate in 5-byte blocks; when finished writing, the caller must
// Close the returned encoder to flush any partially written blocks.
func NewEncoder(w io.Writer) io.WriteCloser {
	return &encoder{w: w}
}
