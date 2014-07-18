# mpeg_go

A pure golang implementation of an MPEG-2 decoder for
educational purposes only.

## Roadmap
- [ ] MPEG-2 Transport Stream (TS) Support
  - [x] Basic parser
  - [ ] Adapation Feild support
  - [x] Demultplexer
- [ ] MPEG-2 Packetized Elementary Stream (PES) Support
  - [x] Basic parser
  - [ ] PES Extension support
- [ ] Program Association Table support
- [ ] Conditional Access Table support
- [ ] Program Map Table support
- [ ] MPEG-2 Program Stream (PS) Support
  - [x] Basic decoder
- [ ] MPEG-2 Video (13818-2) Stream support

## Examples

### Give it a spin!

```
go get github.com/32bitkid/mpeg_go
```


### Demux a TS for a particular PID (0x21)

```go
import "github.com/32bitkid/mpeg_go/ts"

// Implementation of `mpeg_go.BitReader`
import br "github.com/32bitkid/bitreader"

import "os"

func main() {
	file, _ := os.Open("source.ts")
  
	demux := ts.NewDemuxer(br.NewReader32(file))
	packets := demux.Where(ts.IsPID(0x21))
	demux.Go()
	for packet := range packets {
		// Do work!
	}
}
```

### Demux a TS with multiple streams

```go
import "github.com/32bitkid/mpeg_go/ts"


// Implementation of `mpeg_go.BitReader`
import br "github.com/32bitkid/bitreader"

import "os"

func main() {
	file, _ := os.Open("source.ts")
  
	demux := ts.NewDemuxer(br.NewReader32(file))
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

### Build decoder pipelines

```go
import "github.com/32bitkid/mpeg_go/ts"
import "github.com/32bitkid/mpeg_go/pes"

// Implementation of `mpeg_go.BitReader`
import br "github.com/32bitkid/bitreader"

import "os"

func main() {
	file, _ := os.Open("source.ts")
	
	pid := ts.IsPID(0x21)
	
	demux := ts.NewDemuxer(br.NewReader32(file))
	demux.SkipUntil(pid.And(ts.IsPayloadUnitStart))
	
	pesDecoder := pes.NewDecoder()

	// file -> TS packets -> filtered to pid -> PES packets -> ES data
	pesPayload := pesDecoder.TS(demux.Where(pid)).PayloadOnly()
	
	for es := range pesPayload {
	  // do work with Elementary Stream data
	}
}
```
