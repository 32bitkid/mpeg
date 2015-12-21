package pes

import "io"
import "bytes"
import "github.com/32bitkid/bitreader"

func NewPayloadReader(source io.Reader) io.Reader {
	return &payloadReader{
		br:            bitreader.NewBitReader(source),
		currentPacket: new(Packet),
	}
}

type payloadReader struct {
	br            bitreader.BitReader
	currentPacket *Packet
	remainder     bytes.Buffer
}

func (r *payloadReader) Read(p []byte) (n int, err error) {
	// Drain remainder
	for len(p) > 0 {
		cn, err := r.remainder.Read(p)
		n += cn
		p = p[cn:]

		if err == io.EOF {
			break
		} else if err != nil {
			return n, err
		}
	}

	var remainder []byte

	// Fill from packet stream
	for len(p) > 0 {
		err := r.currentPacket.readFrom(r.br)
		if err != nil {
			return n, err
		}

		cn := copy(p, r.currentPacket.Payload)
		n += cn
		p = p[cn:]
		remainder = r.currentPacket.Payload[cn:]
	}

	// Fill remainder
	_, err = r.remainder.Write(remainder)

	return
}
