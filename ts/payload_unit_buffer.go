package ts

import "io"
import "bytes"
import "github.com/32bitkid/mpeg/util"

type PayloadUnitBuffer interface {
	io.Reader
	Fill() (int, error)
}

func NewPayloadUnitBuffer(source io.Reader, where PacketTester) PayloadUnitBuffer {
	return &payloadUnitBuffer{
		br:             util.NewSimpleBitReader(source),
		where:          where,
		startIndicator: where.And(IsPayloadUnitStart),
	}
}

type payloadUnitBuffer struct {
	currentPacket  *Packet
	buffer         bytes.Buffer
	br             util.BitReader32
	where          PacketTester
	startIndicator PacketTester
}

func (stream *payloadUnitBuffer) Read(p []byte) (n int, err error) {
	return stream.buffer.Read(p)
}

func (stream *payloadUnitBuffer) Fill() (n int, err error) {

	stream.init()
	// step until first start indicator
	for {
		isStart := stream.startIndicator(stream.currentPacket)
		if isStart {
			break
		}
		err = stream.advance()
		if err != nil {
			return
		}
	}

	// Read until next start indicator
	for {
		if stream.where(stream.currentPacket) {
			cn, err := stream.buffer.Write(stream.currentPacket.Payload)
			n += cn
			if err != nil {
				break
			}
		}

		err = stream.advance()
		if err != nil {
			break
		}

		if stream.startIndicator(stream.currentPacket) {
			break
		}
	}

	return

}

func (stream *payloadUnitBuffer) init() error {
	// Initialize
	if stream.currentPacket == nil {
		stream.currentPacket = new(Packet)
		err := stream.advance()
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func (stream *payloadUnitBuffer) advance() error {
	return stream.currentPacket.ReadFrom(stream.br)
}
