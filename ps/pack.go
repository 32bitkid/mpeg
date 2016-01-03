package ps

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg/pes"

type Pack struct {
	*PackHeader
	packs chan *pes.Packet
}

func (p *Pack) Packets() pes.PacketChannel {
	return p.packs
}

// TODO better errors?
func readPack(r bitreader.BitReader) (*Pack, <-chan bool, error) {
	var err error

	packComplete := make(chan bool)

	pack := Pack{
		packs: make(chan *pes.Packet),
	}

	pack.PackHeader, err = readPackHeader(r)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(packComplete)
		defer func() { packComplete <- true }()
		defer close(pack.packs)

		for true {

			if nextbits, err := r.Peek32(24); err != nil {
				return
			} else if nextbits != StartCodePrefix {
				break
			}

			if nextbits, err := r.Peek32(32); err != nil {
				return
			} else if StartCode(nextbits) == PackStartCode || StartCode(nextbits) == ProgramEndCode {
				return
			}

			if packet, err := pes.NewPacket(r); err != nil {
				return
			} else {
				pack.packs <- packet
			}
		}
	}()

	return &pack, packComplete, nil
}
