package pes

import "bytes"
import "github.com/32bitkid/mpeg/ts"
import "github.com/32bitkid/bitreader"

type Decoder interface {
	Output() PacketChannel
	Err() error
}

type decoder struct {
	outputs []chan *Packet
	err     error
}

func (d *decoder) Output() PacketChannel {
	newOutput := make(chan *Packet)
	d.outputs = append(d.outputs, newOutput)
	return newOutput
}

func (d *decoder) Err() error {
	return d.err
}

func NewTsDecoder(newFn bitreader.NewReader32Fn, input ts.PacketChannel) Decoder {
	d := decoder{}

	buffer := &bytes.Buffer{}
	reader := newFn(buffer)

	go func() {
		defer func() {
			for _, output := range d.outputs {
				close(output)
			}
		}()

		for tsPacket := range input {

			if tsPacket.PayloadUnitStartIndicator && buffer.Len() > 0 {
				// Drain
				pesPacket, err := ReadPacket(reader, buffer.Len())
				if err != nil {
					d.err = err
					return
				}
				for _, output := range d.outputs {
					output <- pesPacket
				}
			}

			// Fill
			buffer.Write(tsPacket.Payload)
		}
	}()

	return &d
}
