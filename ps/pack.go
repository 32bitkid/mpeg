package ps

import . "github.com/32bitkid/mpeg_go"
import "github.com/32bitkid/mpeg_go/pes"

type Pack struct {
	*PackHeader
	packs chan *pes.Packet
}

func (p *Pack) Packets() pes.PacketChannel {
	return p.packs
}

// TODO better errors?
func readPack(r BitReader) (*Pack, <-chan bool, error) {
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
			v, err := r.Peek32(24)
			if err != nil {
				return
			}

			if v != PacketStartCodePrefix {
				break
			}

			v, err = r.Peek32(32)
			if err != nil {
				return
			}

			if v == PackStartCode || v == ProgramEndCode {
				return
			}

			packet, err := pes.ReadPacket(r, -1)

			if err != nil {
				return
			}

			pack.packs <- packet
		}
	}()

	return &pack, packComplete, nil
}
