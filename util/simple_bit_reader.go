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

func (b *simpleReader32) Peek32(bits uint) (uint32, error) {
	err := b.check(bits)
	if err != nil {
		return 0, err
	}

	shift := (64 - bits)
	var mask uint64 = (1 << (bits + 1)) - 1
	return uint32(b.buffer & (mask << shift) >> shift), err
}

func (b *simpleReader32) Trash(bits uint) error {
	err := b.check(bits)
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	} else if err != nil {
		return err
	}
	b.buffer <<= bits
	b.bitsLeft -= bits
	return err
}

func (b *simpleReader32) Read32(bits uint) (uint32, error) {
	val, err := b.Peek32(bits)
	if err == io.EOF {
		return 0, io.ErrUnexpectedEOF
	} else if err != nil {
		return 0, err
	}
	err = b.Trash(bits)
	return val, err
}

func (b *simpleReader32) PeekBit() (bool, error) {
	val, err := b.Peek32(1)
	return val == 1, err
}

func (b *simpleReader32) ReadBit() (bool, error) {
	val, err := b.PeekBit()
	if err == io.EOF {
		return false, io.ErrUnexpectedEOF
	} else if err != nil {
		return false, err
	}
	err = b.Trash(1)
	return val, err
}

func (b *simpleReader32) check(bits uint) error {
	if b.bitsLeft < bits {
		return b.fill(bits)
	}
	return nil
}

func (b *simpleReader32) fill(needed uint) error {
	neededBytes := int((needed - b.bitsLeft + 7) >> 3)
	n, err := io.ReadAtLeast(b.source, b.readBuffer, neededBytes)

	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
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
