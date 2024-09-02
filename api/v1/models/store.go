package models

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	api "server/api/v1"
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

func (s *store) Read(in int32) (out *api.Record, err error) {
	p := make([]byte, lenWidth)
	n, err := s.ReadAt(p, in)

	if err != nil {
		return nil, err
	}

	if n != lenWidth {
		return nil, io.EOF
	}

	out = &api.Record{
		Value:  p,
		Offset: uint64(in),
	}

	return out, nil
}

func (s *store) ReadAt(p []byte, off int32) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, int64(off))
}

func (s *store) Append(record *api.Record) (bytes int, off uint32, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, 0, err
	}

	off = uint32(s.size)
	if err := binary.Write(s.buf, enc, record.Value); err != nil {
		return 0, 0, err
	}

	if err = binary.Write(s, enc, record.Value); err != nil {
		return 0, 0, err
	}
	s.size += 1

	return len(record.Value), off, nil
}

func (s *store) Remove() error {
	if err := s.Close(); err != nil {
		return err
	}
	return os.Remove(s.Name())
}
