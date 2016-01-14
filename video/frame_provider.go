package video

import "errors"
import "image"

var EOS = errors.New("end of sequence")
var ErrUnsupportedVideoStream_ISO_IEC_11172_2 = errors.New("unsupported video stream ISO/IEC 11172-2")

// Next will return the next frame of video decoded from the video stream.
//
// Please note, this function will return frames in the order they are *decoded*, which may not
// be the order they should be displayed.
func (self *VideoSequence) Next() (image.Image, error) {

	if self.SequenceHeader != nil {
		goto RESUME
	}

	// align to next start code
	if err := next_start_code(self); err != nil {
		return nil, err
	}

	// read sequence_header
	if err := self.sequence_header(); err != nil {
		return nil, err
	}

	// peek for sequence_extension
	if val, err := self.Peek32(32); err != nil {
		return nil, err
	} else if StartCode(val) != ExtensionStartCode {
		// Stream is MPEG-1 Video
		return nil, ErrUnsupportedVideoStream_ISO_IEC_11172_2
	}

	if err := self.sequence_extension(); err != nil {
		return nil, err
	}

CONTINUE:

	if err := self.extension_and_user_data(0); err != nil {
		return nil, err
	}

MORE_FRAMES:

	if nextbits, err := self.Peek32(32); err != nil {
		return nil, err
	} else if StartCode(nextbits) == GroupStartCode {
		if err := self.group_of_pictures_header(); err != nil {
			return nil, err
		}
		self.frameStore.gop()
		if err := self.extension_and_user_data(1); err != nil {
			return nil, err
		}
	}

	if err := self.picture_header(); err != nil {
		return nil, err
	}

	if err := self.picture_coding_extension(); err != nil {
		return nil, err
	}

	if err := self.extension_and_user_data(2); err != nil {
		return nil, err
	}

	self.frameStore.set(self.PictureHeader.temporal_reference)

	if frame, err := self.picture_data(); err != nil {
		return nil, err
	} else {
		switch self.PictureHeader.picture_coding_type {
		case IFrame, PFrame:
			self.frameStore.add(frame, self.PictureHeader.temporal_reference)
		}
		return frame, nil
	}

RESUME:

	if nextbits, err := self.Peek32(32); err != nil {
		return nil, err
	} else if StartCode(nextbits) == PictureStartCode {
		goto MORE_FRAMES
	} else if StartCode(nextbits) == GroupStartCode {
		goto MORE_FRAMES
	}

	if nextbits, err := self.Peek32(32); err != nil {
		return nil, err
	} else if StartCode(nextbits) == SequenceEndStartCode {
		// consume SequenceEndStartCode
		if err := self.Trash(32); err != nil {
			return nil, err
		}
		goto END_OF_STREAM
	}

	if err := self.sequence_header(); err != nil {
		return nil, err
	}

	if err := self.sequence_extension(); err != nil {
		return nil, err
	}

	goto CONTINUE

END_OF_STREAM:
	return nil, EOS
}
