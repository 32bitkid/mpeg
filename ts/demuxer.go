package ts

import "io"

func alwaysTrue(p *Packet) bool { return true }

func Demux(source io.Reader) Demuxer {
	reader := NewReader(source)
	return &tsDemuxer{
		reader:    reader,
		skipUntil: alwaysTrue,
		takeWhile: alwaysTrue,
	}
}

type Demuxer interface {
	Where(PacketTester) <-chan *Packet
	Go() <-chan bool
	Err() error

	SkipUntil(PacketTester)
	TakeWhile(PacketTester)
}

type conditionalChannel struct {
	test    PacketTester
	channel chan<- *Packet
}

type tsDemuxer struct {
	reader             TransportStreamReader
	registeredChannels []conditionalChannel
	lastErr            error
	skipUntil          PacketTester
	takeWhile          PacketTester
}

func (tsd *tsDemuxer) Where(test PacketTester) <-chan *Packet {
	channel := make(chan *Packet)
	tsd.registeredChannels = append(tsd.registeredChannels, conditionalChannel{test, channel})
	return channel
}

func (tsd *tsDemuxer) SkipUntil(skipUntil PacketTester) {
	tsd.skipUntil = skipUntil
}

func (tsd *tsDemuxer) TakeWhile(takeWhile PacketTester) {
	tsd.takeWhile = takeWhile
}

func (tsd *tsDemuxer) Go() <-chan bool {

	done := make(chan bool)
	var skipping = true
	var skipUntil = tsd.skipUntil
	var takeWhile = tsd.takeWhile

	go func() {

		for true {
			p, err := tsd.reader.Next()

			if err != nil {
				tsd.lastErr = err
				done <- true
				break
			}

			if skipping {
				if !skipUntil(p) {
					continue
				} else {
					skipping = false
				}
			} else {
				if !takeWhile(p) {
					break
				}
			}

			for _, item := range tsd.registeredChannels {
				if item.test(p) {
					item.channel <- p
				}
			}

		}
		done <- true
	}()

	return done
}

func (tsd *tsDemuxer) Err() error {
	return tsd.lastErr
}
