package pes

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
