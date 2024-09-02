package models

import (
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth = 4
	posWidth = 4
	entWidth = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}
