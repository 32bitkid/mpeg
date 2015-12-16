package ts

import "bytes"
import "io/ioutil"

type PacketChannel <-chan *Packet

func (input PacketChannel) PayloadOnly() <-chan []byte {
	output := make(chan []byte)
	go func() {
		defer close(output)
		for packet := range input {
			output <- packet.Payload
		}
	}()
	return output
}

func (input PacketChannel) PayloadUnit() <-chan []byte {
	var buf bytes.Buffer
	output := make(chan []byte)
	started := false
	go func() {
		defer close(output)
		for packet := range input {
			buf.Write(packet.Payload)

			if packet.PayloadUnitStartIndicator {
				if started {
					data, err := ioutil.ReadAll(&buf)
					if err != nil {
						output <- data
					}
				} else {
					started = true
				}
			}
		}
	}()
	return output
}
