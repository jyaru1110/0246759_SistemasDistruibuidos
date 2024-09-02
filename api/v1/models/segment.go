package models

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset int64
	config                 Config
}
