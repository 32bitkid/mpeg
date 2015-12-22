package ts

import "io"

// TransportStreamControl is the interface that
// contains functions to limit a transport stream
// packets
type TransportStreamControl interface {
	SkipUntil(skipUntil PacketTester)
	TakeWhile(takeWhile PacketTester)
}

// StreamControlReader is the interface that wraps the basic reading
// function and transport stream control functions.
type StreamControlReader interface {
	io.Reader
	TransportStreamControl
}
