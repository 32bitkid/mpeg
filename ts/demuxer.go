package ts

import "github.com/32bitkid/bitreader"
import "io"

// Creates a new MPEG-2 Transport Stream Demultiplexer
func NewDemuxer(reader io.Reader) Demuxer {
	return &tsDemuxer{
		reader:    bitreader.NewBitReader(reader),
		skipUntil: alwaysTrueTester,
		takeWhile: alwaysTrueTester,
	}
}

// Demuxer is the interface to control and extract
// streams out of a Multiplexed Transport Stream.
type Demuxer interface {
	Where(PacketTester) PacketChannel
	Go() <-chan bool
	Err() error

	TransportStreamControl
}

// Wraps a condition and a channel. Any packets
// that match the PacketTester should be delivered
// to the channel
type conditionalChannel struct {
	tester  PacketTester
	channel chan<- *Packet
}

type tsDemuxer struct {
	reader             bitreader.BitReader
	registeredChannels []conditionalChannel
	lastErr            error
	skipUntil          PacketTester
	takeWhile          PacketTester
}

// Create a Packet Channel that will only include packets
// that match the PacketTester
func (tsd *tsDemuxer) Where(tester PacketTester) PacketChannel {
	channel := make(chan *Packet)
	tsd.registeredChannels = append(tsd.registeredChannels, conditionalChannel{tester, channel})
	return channel
}

// Skip any packets from the input stream until the PacketTester
// returns true
func (tsd *tsDemuxer) SkipUntil(skipUntil PacketTester) {
	tsd.skipUntil = skipUntil
}

// Only return packets from the stream while the PacketTester
// returns true
func (tsd *tsDemuxer) TakeWhile(takeWhile PacketTester) {
	tsd.takeWhile = takeWhile
}

// Create a goroutine to begin parsing the input stream
func (tsd *tsDemuxer) Go() <-chan bool {

	done := make(chan bool)
	var skipping = true
	var skipUntil = tsd.skipUntil
	var takeWhile = tsd.takeWhile
	var p = &Packet{}

	go func() {

		defer func() {
			for _, item := range tsd.registeredChannels {
				close(item.channel)
			}
			done <- true
		}()

		for {
			err := p.ReadFrom(tsd.reader)

			if err != nil {
				tsd.lastErr = err
				return
			}

			if skipping {
				if !skipUntil(p) {
					continue
				} else {
					skipping = false
				}
			} else {
				if !takeWhile(p) {
					return
				}
			}

			for _, item := range tsd.registeredChannels {
				if item.tester(p) {
					item.channel <- p
				}
			}
		}
	}()

	return done
}

// Retrieve the last error from the demuxer
func (tsd *tsDemuxer) Err() error {
	return tsd.lastErr
}
