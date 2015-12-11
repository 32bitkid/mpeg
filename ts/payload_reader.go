package ts

import "io"
import "bytes"
import "github.com/32bitkid/mpeg/util"

type PayloadReader interface {
	io.Reader
	SkipUntil(PacketTester) error
}

func NewPayloadReader(source io.Reader, tester PacketTester) PayloadReader {
	return &payloadReader{
		currentPacket: nil,
		tsr:           util.NewSimpleBitReader(source),
		tester:        tester,
	}
}

type payloadReader struct {
	currentPacket *Packet
	tsr           util.BitReader32
	tester        PacketTester
	remainder     bytes.Buffer
}

func (r *payloadReader) Next() error {
	if r.currentPacket == nil {
		r.currentPacket = new(Packet)
	}
	return r.currentPacket.ReadFrom(r.tsr)
}

func (r *payloadReader) SkipUntil(begin PacketTester) (err error) {
	r.remainder.Reset()
	t := r.tester.And(begin)
	for {
		r.Next()
		if err != nil {
			return err
		}
		if t(r.currentPacket) {
			r.remainder.Write(r.currentPacket.Payload)
			return nil
		}
	}
	return io.EOF
}

func (r *payloadReader) Read(p []byte) (n int, err error) {

	var remainder []byte

	if r.remainder.Len() > 0 {
		copied, err := r.remainder.Read(p)
		if err != nil {
			return copied, err
		}
		n = n + copied
		p = p[copied:]
	}

	for len(p) > 0 {
		err = r.Next()
		if err != nil {
			return
		}

		copied := copy(p, r.currentPacket.Payload)
		n = n + copied
		p = p[copied:]
		remainder = r.currentPacket.Payload[copied:]
	}

	_, err = r.remainder.Write(remainder)

	if err != nil {
		return
	}

	return
}
