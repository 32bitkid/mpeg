package ts

import br "github.com/32bitkid/bitreader"

// Creates a new MPEG-2 Transport Stream Demultiplexer
func NewDemuxer(reader br.Reader32) Demuxer {
	return &tsDemuxer{
		reader:    reader,
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

	SkipUntil(PacketTester) Demuxer
	TakeWhile(PacketTester) Demuxer
}

// Wraps a condition and a channel. Any packets
// that match the PacketTester should be delivered
// to the channel
type conditionalChannel struct {
	test    PacketTester
	channel chan<- *Packet
}

type tsDemuxer struct {
	reader             br.Reader32
	registeredChannels []conditionalChannel
	lastErr            error
	skipUntil          PacketTester
	takeWhile          PacketTester
}

// Create a Packet Channel that will only include Transport Stream
// packets that match the PacketTester
func (tsd *tsDemuxer) Where(test PacketTester) PacketChannel {
	channel := make(chan *Packet)
	tsd.registeredChannels = append(tsd.registeredChannels, conditionalChannel{test, channel})
	return channel
}

// Skip any packets from the input stream until the PacketTester
// returns true
func (tsd *tsDemuxer) SkipUntil(skipUntil PacketTester) Demuxer {
	tsd.skipUntil = skipUntil
	return tsd
}

// Only return packets from the stream while the PacketTester
// returns true
func (tsd *tsDemuxer) TakeWhile(takeWhile PacketTester) Demuxer {
	tsd.takeWhile = takeWhile
	return tsd
}

// Create a goroutine to begin parsing the input stream
func (tsd *tsDemuxer) Go() <-chan bool {

	done := make(chan bool, 1)
	var skipping = true
	var skipUntil = tsd.skipUntil
	var takeWhile = tsd.takeWhile

	go func() {

		defer func() { done <- true }()
		defer func() {
			for _, item := range tsd.registeredChannels {
				close(item.channel)
			}
		}()

		for true {
			p, err := ReadPacket(tsd.reader)

			if err != nil {
				tsd.lastErr = err
				done <- true
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
				if item.test(p) {
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
