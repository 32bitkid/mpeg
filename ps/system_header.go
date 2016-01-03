package ps

import "github.com/32bitkid/bitreader"

type SystemHeader struct {
	HeaderLength          uint32
	RateBound             uint32
	AudioBound            uint32
	Fixed                 bool
	CSPS                  bool
	SystemAudioLock       bool
	SystemVideoLock       bool
	VideoBound            uint32
	PacketRateRestriction bool
	Streams               []*SystemHeaderStream
}

type SystemHeaderStream struct {
	StreamID               uint32
	P_STD_BufferBoundScale bool
	P_STD_BufferSizeBound  uint32
}

func readSystemHeader(r bitreader.BitReader) (*SystemHeader, error) {

	if err := r.Trash(32); err != nil {
		return nil, err
	}

	systemHeader := SystemHeader{}

	if hl, err := r.Read32(16); err != nil {
		return nil, err
	} else {
		systemHeader.HeaderLength = hl
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if rb, err := r.Read32(22); err != nil {
		return nil, err
	} else {
		systemHeader.RateBound = rb
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if ab, err := r.Read32(6); err != nil {
		return nil, err
	} else {
		systemHeader.AudioBound = ab
	}

	if fixed, err := r.ReadBit(); err != nil {
		return nil, err
	} else {
		systemHeader.Fixed = fixed
	}

	if csps, err := r.ReadBit(); err != nil {
		return nil, err
	} else {
		systemHeader.CSPS = csps
	}

	if sal, err := r.ReadBit(); err != nil {
		return nil, err
	} else {
		systemHeader.SystemAudioLock = sal
	}

	if svl, err := r.ReadBit(); err != nil {
		return nil, err
	} else {
		systemHeader.SystemVideoLock = svl
	}

	if v, err := r.Read32(1); err != nil {
		return nil, err
	} else if v != 1 {
		return nil, ErrMarkerNotFound
	}

	if ab, err := r.Read32(5); err != nil {
		return nil, err
	} else {
		systemHeader.AudioBound = ab
	}

	if prr, err := r.ReadBit(); err != nil {
		return nil, err
	} else {
		systemHeader.PacketRateRestriction = prr
	}

	if err := r.Trash(7); err != nil {
		return nil, err
	}

	for true {
		if nextbits, err := r.Peek32(1); err != nil {
			return nil, err
		} else if nextbits != 1 {
			break
		}

		stream := SystemHeaderStream{}

		if sid, err := r.Read32(8); err != nil {
			return nil, err
		} else {
			stream.StreamID = sid
		}

		if v, err := r.Read32(2); err != nil {
			return nil, err
		} else if v != 3 {
			return nil, ErrMarkerNotFound
		}

		if scale, err := r.ReadBit(); err != nil {
			return nil, err
		} else {
			stream.P_STD_BufferBoundScale = scale
		}

		if size, err := r.Read32(13); err != nil {
			return nil, err
		} else {
			stream.P_STD_BufferSizeBound = size
		}

		systemHeader.Streams = append(systemHeader.Streams, &stream)
	}

	return &systemHeader, nil
}
