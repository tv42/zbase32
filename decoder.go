package zbase32

import "io"

type decoder struct {
	io.Reader
	r    io.Reader
	buf  [1024]byte // buffered data waiting to read.
	nbuf int        // the number of bytes in buf
	eof  bool       // indicates that the underlying reader has reached EOF
	err  error
}

func (d *decoder) Read(p []byte) (int, error) {
	var n int

	if d.nbuf < 1 && !d.eof {
		buf := make([]byte, 640)
		l, err := d.r.Read(buf)
		if io.EOF == err {
			d.eof = true
		} else if nil != err {
			return n, err
		}
		if d.nbuf, err = decode(d.buf[0:], buf[:l], -1); nil != err {
			return n, err
		}
	}

	for n < len(p) && d.nbuf > 0 {
		m := copy(p[n:], d.buf[:(min(d.nbuf, len(p)))])
		d.nbuf -= m
		for i := 0; i < d.nbuf; i++ {
			d.buf[i] = d.buf[i+m]
		}
		n += m
	}

	if d.eof == true && d.nbuf == 0 {
		return n, io.EOF
	}

	return n, nil
}

// NewDecoder returns a new z-base-32 stream decoder. Data read from the
// returned reader will be read from r and then decoded.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
