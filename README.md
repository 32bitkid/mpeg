# mpeg

Package mpeg provides an implementation of an  experimental
pure golang implementation of an MPEG-2 decoder. It
is intended as an educational look at some of the patterns and
algorithms involved in the ubiquitous technology of video
compression.


## Roadmap
- [ ] MPEG-2 Transport Stream (TS) Support
  - [x] Basic parser
  - [ ] Adapation Feild support
  - [x] `chan` Demultplexer interface
  - [x] PayloadReader (implementing `io.Reader`)
  - [x] PayloadUnitReader (implementing `io.Reader`)

  - [ ] Other TS Packets
    - [ ] Program Association Table support
    - [ ] Conditional Access Table support
    - [ ] Program Map Table support

- [ ] MPEG-2 Packetized Elementary Stream (PES) Support
  - [x] Basic parser
  - [ ] PES Extension support
  - [ ] `chan` Packet streamer interface
  - [x] PayloadReader (implementing `io.Reader`)

- [ ] MPEG-2 Program Stream (PS) Support
  - [x] Basic decoder
  - [x] `chan` Pack streamer interface
  - [ ] PackReader (implementing `io.Reader`)

- [ ] MPEG-2 Video (13818-2) Stream support
  - [x] I-Frame bitstream decoding
  - [ ] I-Frame renderer

  - [ ] P-Frame bitstream decoding
  - [ ] B-Frame bitstream decoding
  - [ ] Motion vector support
  - [ ] P/B-Frame Renderer


## Examples

### Give it a spin!

```
go get -d github.com/32bitkid/mpeg
```

### Using the `io.Reader` interface

Decode a frame of video and save it as a png

```go
package main

import "os"
import "image/png"

import "github.com/32bitkid/bitreader"
import "github.com/32bitkid/mpeg/ts"
import "github.com/32bitkid/mpeg/pes"
import "github.com/32bitkid/mpeg/video"

func main() {
  tsReader, err := os.Open("source.ts")
  // Decode PID 0x21 from the TS stream
  pesReader := ts.NewPayloadUnitReader(tsReader, ts.IsPID(0x21))
  // Decode the PES stream
  esReader := pes.NewPayloadReader(pesReader)
  // Decode the ES into a stream of frames
  v := video.NewVideoSequence(esReader)

  // Align to next sequence start/entry point
  v.AlignTo(video.SequenceHeaderStartCode)

  // get the next frame
  frame, _ = v.Next()
  file, _ := os.Create("output.png")
  png.Encode(file, frame)
}
```


### Using the streaming interface

#### Demux a TS for a particular PID (0x21)

```go
import "github.com/32bitkid/mpeg/ts"
import "github.com/32bitkid/bitreader"

import "os"

func main() {
	file, _ := os.Open("source.ts")

	demux := ts.NewDemuxer(bitreader.NewBitReader32(file))
	packets := demux.Where(ts.IsPID(0x21))
	demux.Go()
	for packet := range packets {
		// Do work!
	}
}
```

#### Demux a TS with multiple streams

```go
import "github.com/32bitkid/mpeg/ts"
import "github.com/32bitkid/bitreader"

import "os"

func main() {
	file, _ := os.Open("source.ts")

	demux := ts.NewDemuxer(bitreader.NewBitReader32(file))
	hd := demux.Where(ts.IsPID(0x21))
	sd := demux.Where(ts.IsPID(0x31))
	demux.Go()

	var done = false
 	for done == false {
		select {
		case hdPacket := <-hd:
			// process an hd packet
		case sdPacket := <-sd:
			// process an sd packet
		case <-stop:
			done = true
		}
	}
}
```
