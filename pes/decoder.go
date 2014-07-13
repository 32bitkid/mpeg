package pes

import "log"
import "github.com/32bitkid/bitreader"

type PacketChannel <-chan *Packet

func (input PacketChannel) PayloadOnly() <-chan []byte {
	output := make(chan []byte)
	go func() {
		for packet := range input {
			output <- packet.Payload
		}
		close(output)
	}()
	return output
}

func Decoder(input <-chan []byte) PacketChannel {
	output := make(chan *Packet)
	reader := bitreader.NewBufferedBitreader()
	closed := false

	// Sink
	go func() {
		for payload := range input {
			reader.Write(payload)
		}
		closed = true
	}()

	// Pump
	go func() {
		for !closed {
			packet, err := ReadPacket(reader)
			if err != nil {
				log.Println(err)
				close(output)
				return
			}
			log.Printf("%b", packet.StreamID)
		}
	}()

	return output
}
