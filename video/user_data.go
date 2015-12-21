package video

import "github.com/32bitkid/bitreader"

type UserData []byte

func user_data(br bitreader.BitReader) (UserData, error) {
	err := start_code_check(br, UserDataStartCode)
	if err != nil {
		return nil, err
	}

	data := make(UserData, 0)

	for {
		peek, err := br.Peek32(24)
		if err != nil {
			return nil, err
		}

		if peek == StartCodePrefix {
			break
		}

		raw, err := br.Read32(8)
		if err != nil {
			return nil, err
		}

		data = append(data, byte(raw))
	}

	return data, next_start_code(br)

}
