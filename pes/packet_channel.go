package pes

// PacketChannel is a delivery channel of PES Packets
type PacketChannel <-chan *Packet

// PayloadOnly transforms a PacketChannel into a delivery channel of packet payload
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
