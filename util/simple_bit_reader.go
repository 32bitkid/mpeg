package util

import "io"

func NewSimpleBitReader(r io.Reader) BitReader32 {
	return &simpleReader32{r, make([]byte, 4), 0, 0}
}

type simpleReader32 struct {
	source     io.Reader
	readBuffer []byte
	buffer     uint64
	bitsLeft   uint
}

func (b *simpleReader32) Peek32(len uint) (uint32, error) {
	err := b.check(len)
	if err != nil && err != io.EOF {
		return 0, err
	}

	shift := (64 - len)
	var mask uint64 = (1 << (len + 1)) - 1
	return uint32(b.buffer & (mask << shift) >> shift), err
}

func (b *simpleReader32) Trash(len uint) error {
	err := b.check(len)
	if err != nil && err != io.EOF {
		return err
	}
	b.buffer <<= len
	b.bitsLeft -= len
	return err
}

func (b *simpleReader32) Read32(len uint) (uint32, error) {
	val, err := b.Peek32(len)
	if err != nil && err != io.EOF {
		return 0, err
	}
	err = b.Trash(len)
	return val, err
}

func (b *simpleReader32) PeekBit() (bool, error) {
	val, err := b.Peek32(1)
	return val == 1, err
}

func (b *simpleReader32) ReadBit() (bool, error) {
	val, err := b.PeekBit()
	if err != nil && err != io.EOF {
		return val, err
	}
	err = b.Trash(1)
	return val, err
}

func (b *simpleReader32) check(len uint) error {
	if b.bitsLeft < len {
		return b.fill(len)
	}
	return nil
}

func (b *simpleReader32) fill(needed uint) error {
	neededBytes := int((needed - b.bitsLeft + 7) >> 3)
	len, err := io.ReadAtLeast(b.source, b.readBuffer, neededBytes)

	if err != nil && err != io.EOF {
		return err
	}

	if uint(len*8)+b.bitsLeft < needed {
		return ErrNotAvailable
	}

	for i := 0; i < len; i++ {
		b.buffer = b.buffer | uint64(b.readBuffer[i])<<(64-8-b.bitsLeft)
		b.bitsLeft += 8
	}

	return err
}

func (b *simpleReader32) Read(p []byte) (n int, err error) {
	b.ByteAlign()
	bytes := int((b.bitsLeft + 7) >> 3)
	for i := 0; i < bytes; i++ {
		val, err := b.Read32(8)
		if err != nil {
			return i, err
		}
		p[i] = byte(val)
	}
	n, err = b.source.Read(p[bytes:])
	n += bytes
	return
}

func (b *simpleReader32) ByteAlign() error {
	return b.Trash(b.bitsLeft % 8)
}

func (b *simpleReader32) IsByteAligned() bool {
	return b.bitsLeft%8 == 0
}
