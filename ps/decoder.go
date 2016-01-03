package ps

import "github.com/32bitkid/bitreader"

type Decoder interface {
	Packs() <-chan *Pack
	Go() <-chan bool
	Err() error
}

type decoder struct {
	r     bitreader.BitReader
	packs chan *Pack
	err   error
}

func (d *decoder) Go() <-chan bool {
	done := make(chan bool)

	go func() {
		defer close(done)
		defer func() { done <- d.err == nil }()
		defer close(d.packs)

		for true {
			pack, packComplete, err := readPack(d.r)
			if err != nil {
				d.err = err
				return
			}

			d.packs <- pack

			<-packComplete

			if nextbits, err := d.r.Peek32(32); err != nil {
				d.err = err
				return
			} else if StartCode(nextbits) != PackStartCode {
				break
			}
		}

		if nextbits, err := d.r.Peek32(32); err != nil {
			d.err = err
			return
		} else if StartCode(nextbits) != ProgramEndCode {
			d.err = ErrProgramEndCodeNotFound
			return
		}
	}()

	return done
}

func (d *decoder) Packs() <-chan *Pack {
	return d.packs
}

func (d *decoder) Err() error {
	return d.err
}

func NewDecoder(r bitreader.BitReader) Decoder {
	return &decoder{
		r:     r,
		packs: make(chan *Pack),
		err:   nil,
	}
}
