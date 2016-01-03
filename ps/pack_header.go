package ps

import "github.com/32bitkid/bitreader"
import "errors"

type PackHeader struct {
	SystemClockReferenceBase      uint32
	SystemClockReferenceExtension uint32
	ProgramMuxRate                uint32
	*SystemHeader
}

var ErrNoPackStartCode = errors.New("no pack start code")

func readPackHeader(r bitreader.BitReader) (*PackHeader, error) {

	if nextbits, err := r.Read32(32); err != nil {
		return nil, err
	} else if StartCode(nextbits) != PackStartCode {
		return nil, ErrNoPackStartCode
	}

	if nextbits, err := r.Read32(2); err != nil {
		return nil, err
	} else if nextbits != 1 {
		return nil, ErrMarkerNotFound
	}

	packHeader := PackHeader{}

	if v, err := r.Read32(3); err != nil {
		return nil, err
	} else {
		packHeader.SystemClockReferenceBase |= v << 30
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if v, err := r.Read32(15); err != nil {
		return nil, err
	} else {
		packHeader.SystemClockReferenceBase |= v << 15
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if v, err := r.Read32(15); err != nil {
		return nil, err
	} else {
		packHeader.SystemClockReferenceBase |= v
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if scre, err := r.Read32(9); err != nil {
		return nil, err
	} else {
		packHeader.SystemClockReferenceExtension = scre
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if pmr, err := r.Read32(22); err != nil {
		return nil, err
	} else {
		packHeader.ProgramMuxRate = pmr
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if err := r.Trash(5); err != nil {
		return nil, err
	}

	if pack_stuffing_length, err := r.Read32(3); err != nil {
		return nil, err
	} else {
		for pack_stuffing_length > 0 {
			r.Trash(8) // stuffing_byte
			pack_stuffing_length--
		}
	}

	if nextbits, err := r.Peek32(32); err != nil {
		return nil, err
	} else if StartCode(nextbits) == SystemHeaderStartCode {
		packHeader.SystemHeader, err = readSystemHeader(r)
		if err != nil {
			return nil, err
		}
	}

	return &packHeader, nil
}
