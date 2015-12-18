package ps

import "github.com/32bitkid/mpeg/util"

type PackHeader struct {
	SystemClockReferenceBase      uint32
	SystemClockReferenceExtension uint32
	ProgramMuxRate                uint32
	*SystemHeader
}

func readPackHeader(r util.BitReader32) (*PackHeader, error) {
	var (
		v   uint32
		err error
	)

	if v, err := r.Read32(32); v != PackStartCode || err != nil {
		return nil, err
	}

	if v, err := r.Read32(2); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	packHeader := PackHeader{}

	if v, err = r.Read32(3); err != nil {
		return nil, err
	}

	packHeader.SystemClockReferenceBase |= v << 30

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if v, err = r.Read32(15); err != nil {
		return nil, err
	}

	packHeader.SystemClockReferenceBase |= v << 15

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if v, err = r.Read32(15); err != nil {
		return nil, err
	}

	packHeader.SystemClockReferenceBase |= v

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if packHeader.SystemClockReferenceExtension, err = r.Read32(9); err != nil {
		return nil, err
	}

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if packHeader.ProgramMuxRate, err = r.Read32(22); err != nil {
		return nil, err
	}

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if err = r.Trash(5); err != nil {
		return nil, err
	}

	if v, err = r.Read32(3); err != nil {
		return nil, err
	}

	for v > 0 {
		r.Trash(8)
		v--
	}

	v, err = r.Peek32(32)
	if err != nil {
		return nil, err
	}

	if v == SystemHeaderStartCode {
		packHeader.SystemHeader, err = readSystemHeader(r)
		if err != nil {
			return nil, err
		}
	}

	return &packHeader, nil
}
