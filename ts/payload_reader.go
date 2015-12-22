package ts

import "io"
import "bytes"
import "github.com/32bitkid/bitreader"

// NewPayloadReader takes a transport stream and creates a reader
// that delivers just the packet payload bytes.
func NewPayloadReader(source io.Reader, where PacketTester) StreamControlReader {
	return &payloadReader{
		br:        bitreader.NewBitReader(source),
		where:     where,
		closed:    false,
		skipUntil: nil,
		takeWhile: alwaysTrueTester,
	}
}

type payloadReader struct {
	currentPacket *Packet
	br            bitreader.BitReader
	where         PacketTester
	skipUntil     PacketTester
	takeWhile     PacketTester
	remainder     bytes.Buffer
	closed        bool
}

func (r *payloadReader) SkipUntil(skipUntil PacketTester) {
	r.skipUntil = skipUntil
}

func (r *payloadReader) TakeWhile(takeWhile PacketTester) {
	r.takeWhile = takeWhile
}

func (r *payloadReader) Read(p []byte) (n int, err error) {

	if r.closed == true {
		return 0, io.EOF
	}

	if r.currentPacket == nil {
		r.currentPacket = new(Packet)
		if r.skipUntil != nil {
			err = r.realign()
			if err != nil {
				return
			}

		}
	}

	var remainder []byte

	// Drain remainder
	for len(p) > 0 {
		cn, err := r.remainder.Read(p)
		n = n + cn
		p = p[cn:]
		if err == io.EOF {
			break
		} else if err != nil {
			return n, err
		}
	}

	// Fill from packet stream
	for len(p) > 0 {
		err = r.next()
		if err != nil {
			return
		}

		if r.where(r.currentPacket) {
			copied := copy(p, r.currentPacket.Payload)
			n = n + copied
			p = p[copied:]
			remainder = r.currentPacket.Payload[copied:]
		}

		cont := r.takeWhile(r.currentPacket)
		if cont == false {
			r.closed = true
			return n, io.EOF
		}
	}

	_, err = r.remainder.Write(remainder)
	return
}

func (r *payloadReader) next() error {
	return r.currentPacket.ReadFrom(r.br)
}

func (r *payloadReader) realign() (err error) {
	for {
		r.next()
		if err != nil {
			return err
		}
		done := r.skipUntil(r.currentPacket)
		if done {
			r.remainder.Write(r.currentPacket.Payload)
			return nil
		}
	}
}
