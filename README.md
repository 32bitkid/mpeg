# Mpeg-Go

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
go get github.com/32bitkid/mpeg-go
```


### Demux a TS for a particular PID (0x21)

```go
import "github.com/32bitkid/mpeg-go/ts"
import "os"

func main() {
	file, _ := os.Open("source.ts")
  
	demux := ts.NewDemuxer(file)
	packets := demux.Where(ts.IsPID(0x21))
	demux.Go()
	for packet := range packets {
		// Do work!
	}
}
```

### Demux a TS with multiple streams

```go
import "github.com/32bitkid/mpeg-go/ts"
import "os"

func main() {
	file, _ := os.Open("source.ts")
  
	demux := ts.NewDemuxer(file)
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
import "github.com/32bitkid/mpeg-go/ts"
import "github.com/32bitkid/mpeg-go/pes"
import "os"

func main() {
	file, _ := os.Open("source.ts")
	
	pid := ts.IsPID(0x21)
	
	demux := ts.NewDemuxer(file)
	demux.SkipUntil(pid.And(ts.IsPayloadUnitStart))
	
	// file -> TS packets -> filtered to pid -> PES packets -> ES data
	pesPayload := pes.TsDecoder(demux.Where(pid)).PayloadOnly()
	
	stop := demux.Go()
	var done = false
	for done == false {
		select {
		case es := <-hdVideo:
			// es is a []byte
		case <-stop:
			log.Println("End of stream")
			done = true
	}
}
```
