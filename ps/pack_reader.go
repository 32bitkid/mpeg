package ps

import (
	"bytes"
	"io"
)

import (
	"github.com/32bitkid/bitreader"
	"github.com/32bitkid/mpeg/pes"
)

type packreader struct {
	br        bitreader.BitReader
	remainder *bytes.Buffer
}

func NewPackReader(r io.Reader) io.Reader {
	return &packreader{
		br:        bitreader.NewReader(r),
		remainder: &bytes.Buffer{},
	}
}

func (pr packreader) Read(d []byte) (int, error) {
	n := 0

	// DRAIN buffer
	{
		read, err := pr.remainder.Read(d)
		if err != nil && err != io.EOF {
			return read, err
		} else if err != io.EOF {
			n += read
			d = d[read:]
			if len(d) == 0 {
				return n, nil
			}
		}
	}

STEP:
	if nextbits, err := pr.br.Peek32(32); err != nil {
		return 0, nil
	} else if StartCode(nextbits) == ProgramEndCode {
		goto PROGRAM_END
	} else if StartCode(nextbits) == PackStartCode {
		goto PACK
	} else {
		goto MORE_PACKETS
	}

PACK:
	if _, err := NewPackHeader(pr.br); err != nil {
		return 0, err
	}

MORE_PACKETS:

	if p, err := pes.NewPacket(pr.br); err != nil {
		return 0, nil
	} else {
		copied := copy(d, p.Payload)
		n += copied
		d = d[copied:]
		if len(d) == 0 {
			// Re-fill remainder
			if _, err := pr.remainder.Write(p.Payload[copied:]); err != nil {
				return n, err
			}
			return n, nil
		}
	}

	goto STEP

PROGRAM_END:

	pr.br.Skip(32)
	return n, io.EOF
}
