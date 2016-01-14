package video

import "github.com/32bitkid/bitreader"

type GroupOfPicturesHeader struct {
	timeCode   uint32 // 25 bslbf
	ClosedGOP  bool   // 1 uimsbf
	BrokenLink bool   // 1 uimsbf
}

// ReadGOPHeader parses a group_of_pictures header from the given bitstream.
func ReadGOPHeader(br bitreader.BitReader) (*GroupOfPicturesHeader, error) {
	if err := GroupStartCode.Assert(br); err != nil {
		return nil, err
	}

	goph := GroupOfPicturesHeader{}

	if time_code, err := br.Read32(25); err != nil {
		return nil, err
	} else {
		goph.timeCode = time_code
	}

	if closed_gop, err := br.ReadBit(); err != nil {
		return nil, err
	} else {
		goph.ClosedGOP = closed_gop
	}

	if broken_link, err := br.ReadBit(); err != nil {
		return nil, err
	} else {
		goph.BrokenLink = broken_link
	}

	return &goph, next_start_code(br)
}

// TimeCode represents the associated time code with the first picture following the
// Group of Pictures header with a TemporalReference = 0. DropFrame will only be true
// if the desired framerate is 29.97Hz.
//
// TimeCode appears in the bitstream as a 25-bit integer that has the following layout:
//
//  ├1b┤
//  ┌──┬──────────────────┬──────────────────────┬──┬──────────────────────┬──────────────────────┐
//  │DF│Hours             │Minutes               │MB│Seconds               │Pictures              │
//  └──┴──────────────────┴──────────────────────┴──┴──────────────────────┴──────────────────────┘
//  ├───────────────────────────────────────────25 bits───────────────────────────────────────────┤
type TimeCode struct {
	DropFrame bool
	Hours     int
	Minutes   int
	Seconds   int
	Pictures  int
}

// Returns a parsed TimeCode from the raw GOP header data.
func (gop *GroupOfPicturesHeader) TimeCode() TimeCode {
	return TimeCode{
		DropFrame: (gop.timeCode >> 24) == 1,
		Hours:     (gop.timeCode >> 19) & 0x1F,
		Minutes:   (gop.timeCode >> 13) & 0x3F,
		Seconds:   (gop.timeCode >> 6) & 0x3F,
		Pictures:  (gop.timeCode) & 0x3F,
	}
}
