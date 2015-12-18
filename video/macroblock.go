package video

import "github.com/32bitkid/mpeg/util"

type Macroblock struct {
	macroblock_address_increment uint32
	macroblock_type              *MacroblockType
}

func macroblock(br util.BitReader32) (*Macroblock, error) {

	mb := Macroblock{}

	for {
		nextbits, err := br.Peek32(11)
		if err != nil {
			return nil, err
		}
		if nextbits == 0x08 { // 0000 0001 000
			br.Trash(11)
			mb.macroblock_address_increment += 33
		}

		incr, err := MacroblockAddressIncrementDecoder.Decode(br)
		if err != nil {
			return nil, err
		}
		mb.macroblock_address_increment += incr

		panic("not implemented: macroblock")
	}
}
