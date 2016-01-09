package video

type UserData []byte

func (br *VideoSequence) user_data() (UserData, error) {
	if err := UserDataStartCode.assert(br); err != nil {
		return nil, err
	}

	data := make(UserData, 0)

	for {
		if peek, err := br.Peek32(24); err != nil {
			return nil, err
		} else if peek == StartCodePrefix {
			break
		}

		if raw, err := br.Read32(8); err != nil {
			return nil, err
		} else {
			data = append(data, byte(raw))
		}
	}

	return data, next_start_code(br)
}
