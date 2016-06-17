package test

import (
	"encoding/hex"
	"io"
)

type MemDisk struct {
	data       []byte
	woff, roff int
}

func NewMemDisk(n int) *MemDisk {
	return &MemDisk{data: make([]byte, n)}
}

func (w *MemDisk) Write(b []byte) (int, error) {
	n, err := w.WriteAt(b, int64(w.woff))
	w.woff += n
	return n, err
}

func (w *MemDisk) WriteAt(b []byte, off int64) (int, error) {
	for cap(w.data) < len(b)+int(off) {
		w.data = w.data[:cap(w.data)]
		w.data = append(w.data, 0)[:len(w.data)]
	}
	w.data = w.data[:int(off)+len(b)]
	copy(w.data[off:], b)
	return len(b), nil
}

func (w *MemDisk) ReadAt(b []byte, off int64) (int, error) {
	if len(w.data) <= int(off) {
		return 0, io.EOF
	}
	n := copy(b, w.data[off:])
	return n, nil
}

func (w *MemDisk) Dump() string {
	return hex.Dump(w.data)
}

func (w *MemDisk) SeekRead(offset int64, whence int) (ret int64) {
	switch whence {
	case 0:
		w.roff += int(offset)
	case 1:
		w.roff = int(offset)
	default:
	}
	return int64(w.roff)
}

func (w *MemDisk) SeekWrite(offset int64, whence int) (ret int64) {
	switch whence {
	case 0:
		w.woff += int(offset)
	case 1:
		w.woff = int(offset)
	default:
	}
	return int64(w.woff)
}

func (w *MemDisk) Read(b []byte) (int, error) {
	n, err := w.ReadAt(b, int64(w.roff))
	w.roff += n
	return n, err
}

func (w *MemDisk) Close() error {
	return nil
}
