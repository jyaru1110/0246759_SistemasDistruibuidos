package models

import (
	"bufio"
	"fmt"
	"os"
	"path"
	api "server/api/v1"

	"github.com/tysonmote/gommap"
)

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset uint64
	config                 Config
}

func newStore(f *os.File) (*store, error) {
	return &store{
		File: f,
		buf:  bufio.NewWriter(f),
		size: 0,
	}, nil
}

func newIndex(f *os.File, c Config) (*index, error) {
	fileInfo, err := f.Stat()

	if err != nil {
		return nil, err
	}

	f.Truncate(int64(c.Segment.MaxIndexBytes))

	mmap, err := gommap.Map(f.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return &index{
		file: f,
		size: uint64(fileInfo.Size()),
		mmap: mmap,
	}, nil
}

func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}
	var err error
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = newIndex(indexFile, c); err != nil {
		return nil, err
	}
	if off, _, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1
	}

	return s, nil
}

func (s *segment) Append(record *api.Record) (uint64, error) {
	_, pos, err := s.store.Append(record)
	if err != nil {
		return 0, err
	}
	off, err := s.index.Write(record.Offset, pos)
	if err != nil {
		return 0, err
	}
	return uint64(off), nil
}

func (s *segment) Read(off uint64) (*api.Record, error) {
	pos, _, err := s.index.Read(int64(off - s.baseOffset))
	if err != nil {
		return nil, err
	}
	return s.store.Read(int32(pos))
}

func (s *segment) IsMaxed() bool {
	return s.store.size >= s.config.Segment.MaxStoreBytes || s.index.size >= s.config.Segment.MaxIndexBytes
}

func (s *segment) Close() error {
	if err := s.index.Close(); err != nil {
		return err
	}
	if err := s.store.Close(); err != nil {
		return err
	}
	return nil
}

func (s *segment) Remove() error {
	if err := s.index.Remove(); err != nil {
		return err
	}
	if err := s.store.Remove(); err != nil {
		return err
	}
	return nil
}

func (s *segment) Name() string {
	return fmt.Sprintf("%d-%d", s.baseOffset, s.nextOffset)
}