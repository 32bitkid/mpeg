package ps

import "github.com/32bitkid/mpeg/util"

type Decoder interface {
	Packs() <-chan *Pack
	Go() <-chan bool
	Err() error
}

type decoder struct {
	r     util.BitReader32
	packs chan *Pack
	err   error
}

func (d *decoder) Go() <-chan bool {
	done := make(chan bool)

	go func() {
		defer close(done)
		defer func() { done <- d.err == nil }()
		defer close(d.packs)

		var (
			v   uint32
			err error
		)

		for true {

			pack, packComplete, err := readPack(d.r)
			if err != nil {
				d.err = err
				return
			}

			d.packs <- pack

			<-packComplete

			v, err = d.r.Peek32(32)
			if err != nil {
				d.err = err
				return
			}
			if v != PackStartCode {
				break
			}
		}

		v, err = d.r.Peek32(32)
		if v != ProgramEndCode || err != nil {
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

func NewDecoder(r util.BitReader32) Decoder {
	return &decoder{
		r:     r,
		packs: make(chan *Pack),
		err:   nil,
	}
}
