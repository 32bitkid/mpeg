package ts

import "io"
import "github.com/32bitkid/bitreader"

func alwaysTrue(p *Packet) bool { return true }

type PacketChannel <-chan *Packet

func (input PacketChannel) PayloadOnly() <-chan []byte {
	output := make(chan []byte)
	go func() {
		for packet := range input {
			output <- packet.Payload
		}
		close(output)
	}()
	return output
}

func Demux(source io.Reader) Demuxer {
	reader := bitreader.NewReader32(source)
	return &tsDemuxer{
		reader:    reader,
		skipUntil: alwaysTrue,
		takeWhile: alwaysTrue,
	}
}

type Demuxer interface {
	Where(PacketTester) PacketChannel
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
	reader             bitreader.Reader32
	registeredChannels []conditionalChannel
	lastErr            error
	skipUntil          PacketTester
	takeWhile          PacketTester
}

func (tsd *tsDemuxer) Where(test PacketTester) PacketChannel {
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
			p, err := ReadPacket(tsd.reader)

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

		for _, item := range tsd.registeredChannels {
			close(item.channel)
		}

		done <- true
	}()

	return done
}

func (tsd *tsDemuxer) Err() error {
	return tsd.lastErr
}
