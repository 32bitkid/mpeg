package pes

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
