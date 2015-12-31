package ts_test

import (
	"encoding/base64"
	"fmt"
	"github.com/32bitkid/bitreader"
	"github.com/32bitkid/mpeg/ts"
	"strings"
)

// Step through a bit stream one packet at a time.
func Example() {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(shortTsStream))
	br := bitreader.NewBitReader(reader)

	packet := new(ts.Packet)
	for {
		// Read the next packet
		err := packet.Next(br)
		if err != nil {
			break
		}
		fmt.Println(packet)
	}

	// Output:
	// { PID: 0x21, Counter: d }
	// { PID: 0x31, Counter: 6 }
	// { PID: 0x21, Counter: e }
	// { PID: 0x41, Counter: b }
	// { PID: 0x21, Counter: f }
}
