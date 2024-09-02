package models

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func (s *store) Read(in uint64) (out []byte, err error) {
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	value_size_bytes := make([]byte, lenWidth)

	if _, err := s.File.ReadAt(value_size_bytes, int64(in)); err != nil {
		return nil, err
	}

	value_size := enc.Uint64(value_size_bytes)

	value := make([]byte, value_size)

	if _, err := s.File.ReadAt(value, int64(in+lenWidth)); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, int64(off))
}

func (s *store) Append(value []byte) (bytes uint64, off uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, 0, err
	}

	off = s.size
	if err := binary.Write(s.buf, enc, uint64(len(value))); err != nil {
		return 0, 0, err
	}
	if err := binary.Write(s.buf, enc, value); err != nil {
		return 0, 0, err
	}

	s.size += lenWidth + uint64(len(value))

	return uint64(lenWidth) + uint64(len(value)), off, nil
}

func (s *store) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	return os.Remove(s.Name())
}

func newStore(f *os.File) (*store, error) {
	file_info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &store{
		File: f,
		buf:  bufio.NewWriter(f),
		size: uint64(file_info.Size()),
	}, nil
}

func (s *store) Close() error {
	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}
