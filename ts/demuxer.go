package ts

import "io"

func Demux(source io.Reader) Demuxer {
	reader := NewReader(source)
	return &tsDemuxer{
		reader: reader,
	}
}

type PacketPicker func(*TsPacket) bool

type Demuxer interface {
	Where(PacketPicker) <-chan *TsPacket
	PID(uint32) <-chan *TsPacket
	Begin() <-chan bool
	Err() error
}

type packetPickerChannel struct {
	picker  PacketPicker
	channel chan<- *TsPacket
}

type tsDemuxer struct {
	reader         TransportStreamReader
	packetChannels []packetPickerChannel
	lastErr        error
}

func (tsd *tsDemuxer) PID(PID uint32) <-chan *TsPacket {
	return tsd.Where(func(p *TsPacket) bool { return p.PID == PID })
}

func (tsd *tsDemuxer) Where(picker PacketPicker) <-chan *TsPacket {
	channel := make(chan *TsPacket)
	tsd.packetChannels = append(tsd.packetChannels, packetPickerChannel{picker, channel})
	return channel
}

func (tsd *tsDemuxer) Begin() <-chan bool {
	done := make(chan bool)

	go func() {
		for true {
			p, err := tsd.reader.Next()

			if err != nil {
				tsd.lastErr = err
				done <- true
				break
			}

			for _, item := range tsd.packetChannels {
				if item.picker(p) {
					item.channel <- p
				}
			}
		}
	}()

	return done
}

func (tsd *tsDemuxer) Err() error {
	return tsd.lastErr
}
