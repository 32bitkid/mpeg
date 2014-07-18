package pes

import "bytes"
import "github.com/32bitkid/mpeg_go/ts"
import "github.com/32bitkid/bitreader"

type Decoder interface {
	TS(ts.PacketChannel) PacketChannel
	Err() error
}

func NewDecoder() Decoder {
	return &decoder{}
}

type decoder struct {
	err error
}

func (d *decoder) Err() error {
	return d.err
}

func (d *decoder) TS(input ts.PacketChannel) PacketChannel {
	output := make(chan *Packet)

	buffer := &bytes.Buffer{}
	reader := bitreader.NewReader32(buffer)

	go func() {
		defer close(output)

		for tsPacket := range input {

			if tsPacket.PayloadUnitStartIndicator && buffer.Len() > 0 {
				// Drain
				pesPacket, err := ReadPacket(reader, buffer.Len())
				if err != nil {
					d.err = err
					return
				}
				output <- pesPacket
			}

			// Fill
			buffer.Write(tsPacket.Payload)
		}
	}()

	return output
}
