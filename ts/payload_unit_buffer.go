package ts

import "io"
import "bytes"
import "errors"
import "github.com/32bitkid/mpeg/util"

var (
	EOP = errors.New("end of packet")
)

type streamState int

const (
	_                 = iota
	ready streamState = iota
	drained
	refill
)

func NewPayloadUnitBuffer(source io.Reader, where PacketTester) io.Reader {
	return &payloadUnitBuffer{
		br:             util.NewSimpleBitReader(source),
		where:          where,
		startIndicator: where.And(IsPayloadUnitStart),
		state:          refill,
	}
}

type payloadUnitBuffer struct {
	currentPacket  *Packet
	buffer         bytes.Buffer
	br             util.BitReader32
	where          PacketTester
	startIndicator PacketTester
	state          streamState
}

func (stream *payloadUnitBuffer) Read(p []byte) (n int, err error) {
	if stream.state == drained {
		stream.state = refill
		return 0, EOP
	} else if stream.state == refill {
		_, ferr := stream.fill()
		if ferr != nil {
			return 0, ferr
		}
		stream.state = ready
	}

	n, err = stream.buffer.Read(p)

	if err == io.EOF {
		stream.state = drained
		err = nil
	}

	return
}

func (stream *payloadUnitBuffer) fill() (n int, err error) {

	// initialize
	if stream.currentPacket == nil {
		stream.currentPacket = new(Packet)
		err := stream.advance()
		if err != nil && err != io.EOF {
			return 0, err
		}
		err = nil
	}

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

func (stream *payloadUnitBuffer) advance() error {
	return stream.currentPacket.ReadFrom(stream.br)
}
