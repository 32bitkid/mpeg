package ts

import "io"
import "bytes"
import "github.com/32bitkid/mpeg/util"

type PayloadReader interface {
	io.Reader
}

func NewPayloadReader(source io.Reader, tester PacketTester) PayloadReader {
	return &payloadReader{
		packet: new(Packet),
		tsr:    util.NewSimpleBitReader(source),
		tester: tester,
	}
}

type payloadReader struct {
	packet    *Packet
	tsr       util.BitReader32
	tester    PacketTester
	remainder bytes.Buffer
}

func (r *payloadReader) Read(p []byte) (n int, err error) {

	var remainder []byte
	packet := r.packet

	if r.remainder.Len() > 0 {
		copied, err := r.remainder.Read(p)
		if err != nil {
			return copied, err
		}
		n = n + copied
		p = p[copied:]
	}

	for len(p) > 0 {
		err = packet.ReadFrom(r.tsr)
		if err != nil {
			return
		}

		copied := copy(p, packet.Payload)
		n = n + copied
		p = p[copied:]
		remainder = packet.Payload[copied:]
	}

	_, err = r.remainder.Write(remainder)

	if err != nil {
		return
	}

	return
}
