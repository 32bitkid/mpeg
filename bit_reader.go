package mpeg

import "io"

type BitReader interface {
	Read32(uint) (uint32, error)
	Peek32(uint) (uint32, error)
	ReadBit() (bool, error)
	PeekBit() (bool, error)
	Trash(uint) error
	io.Reader
}