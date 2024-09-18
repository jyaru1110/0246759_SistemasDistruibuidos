package models

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4 // 4 bytes for uint64
	posWidth uint64 = 8 // 8 bytes for uint32
	entWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func (i *index) Read(in int64) (uint32, uint64, error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if in == -1 {
		in = int64(i.size / entWidth)
		in--
	}

	pos := in * int64(entWidth)
	if uint64(pos) >= i.size {
		return 0, 0, io.EOF
	}
	off := enc.Uint32(i.mmap[pos : pos+int64(offWidth)])
	idx := enc.Uint64(i.mmap[pos+int64(offWidth) : pos+int64(entWidth)])

	return off, idx, nil
}

func (i *index) Write(offset uint32, pos uint64) error {
	if i.size+uint64(entWidth) > uint64(len(i.mmap)) {
		return io.EOF
	}

	binary.BigEndian.PutUint32(i.mmap[i.size:], offset)

	binary.BigEndian.PutUint64(i.mmap[i.size+uint64(offWidth):], pos)

	i.size += uint64(entWidth)

	return nil
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	i.file.Truncate(int64(i.size))
	if err := i.file.Close(); err != nil {
		return err
	}
	return nil
}

func (i *index) Remove() error {
	if err := os.Remove(i.file.Name()); err != nil {
		return err
	}
	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}
