package video

import br "github.com/32bitkid/bitreader"

func next_start_code(br br.Reader32) error {
	for !br.IsByteAligned() {
		if err := br.Trash(1); err != nil {
			return err
		}
	}

	var err error

	for true {
		if val, err := br.Peek32(24); val == 0x000001 || err != nil {
			break
		}
		if val, err := br.Read32(8); val != 0 || err != nil {
			break
		}
	}

	return err
}
