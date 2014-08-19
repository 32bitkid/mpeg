package pes

import "bytes"
import "github.com/32bitkid/mpeg/ts"
import "github.com/32bitkid/bitreader"

type Decoder interface {
	Output() PacketChannel
	Err() error
}

type decoder struct {
	output chan *Packet
	err    error
}

func (d *decoder) Output() PacketChannel {
	return d.output
}

func (d *decoder) Err() error {
	return d.err
}

func NewTsDecoder(newFn bitreader.NewReader32Fn, input ts.PacketChannel) Decoder {
	d := decoder{
		output: make(chan *Packet),
	}

	buffer := &bytes.Buffer{}
	reader := newFn(buffer)

	go func() {
		defer close(d.output)

		for tsPacket := range input {

			if tsPacket.PayloadUnitStartIndicator && buffer.Len() > 0 {
				// Drain
				pesPacket, err := ReadPacket(reader, buffer.Len())
				if err != nil {
					d.err = err
					return
				}
				d.output <- pesPacket
			}

			// Fill
			buffer.Write(tsPacket.Payload)
		}
	}()

	return &d
}
