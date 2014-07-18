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

func readSystemHeader(r bitreader.Reader32) (*SystemHeader, error) {
	var (
		v   uint32
		err error
	)

	if err = r.Trash(32); err != nil {
		return nil, err
	}

	systemHeader := SystemHeader{}

	if systemHeader.HeaderLength, err = r.Read32(16); err != nil {
		return nil, err
	}

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if systemHeader.RateBound, err = r.Read32(22); err != nil {
		return nil, err
	}

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if systemHeader.AudioBound, err = r.Read32(6); err != nil {
		return nil, err
	}

	if systemHeader.Fixed, err = r.ReadBit(); err != nil {
		return nil, err
	}

	if systemHeader.CSPS, err = r.ReadBit(); err != nil {
		return nil, err
	}

	if systemHeader.SystemAudioLock, err = r.ReadBit(); err != nil {
		return nil, err
	}

	if systemHeader.SystemVideoLock, err = r.ReadBit(); err != nil {
		return nil, err
	}

	if v, err = r.Read32(1); v != 1 || err != nil {
		return nil, ErrMarkerNotFound
	}

	if systemHeader.AudioBound, err = r.Read32(5); err != nil {
		return nil, err
	}

	if systemHeader.PacketRateRestriction, err = r.ReadBit(); err != nil {
		return nil, err
	}

	if err = r.Trash(7); err != nil {
		return nil, err
	}

	for true {
		v, err = r.Peek32(1)

		if err != nil {
			return nil, err
		}

		if v != 1 {
			break
		}

		stream := SystemHeaderStream{}

		if stream.StreamID, err = r.Read32(8); err != nil {
			return nil, err
		}

		if v, err = r.Read32(2); v != 3 || err != nil {
			return nil, err
		}

		if stream.P_STD_BufferBoundScale, err = r.ReadBit(); err != nil {
			return nil, err
		}

		if stream.P_STD_BufferSizeBound, err = r.Read32(13); err != nil {
			return nil, err
		}

		systemHeader.Streams = append(systemHeader.Streams, &stream)
	}

	return &systemHeader, nil
}
