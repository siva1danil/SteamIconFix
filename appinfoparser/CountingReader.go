package appinfoparser

import (
	"errors"
	"io"
	"sync"
)

type CountingReader struct {
	r   io.Reader
	pos int64
	mu  sync.Mutex
}

func (cr *CountingReader) Read(p []byte) (n int, err error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	n, err = cr.r.Read(p)
	cr.pos += int64(n)
	return n, err
}

func (cr *CountingReader) SeekRelative(offset int64) error {
	seeker, ok := cr.r.(io.Seeker)
	if !ok {
		return errors.New("underlying reader does not support seeking")
	}
	_, err := seeker.Seek(offset, io.SeekCurrent)
	if err == nil {
		cr.pos += offset
	}
	return err
}
