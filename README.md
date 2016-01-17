# mpeg

Package mpeg provides an implementation of an  experimental
pure golang implementation of an MPEG-2 decoder. It
is intended as an educational look at some of the patterns and
algorithms involved in the ubiquitous technology of video
compression/decompression.

This package is experimental and is not intended for
use in production environments.

[![GoDoc](https://godoc.org/github.com/32bitkid/mpeg?status.svg)](https://godoc.org/github.com/32bitkid/mpeg)

## Composition

This library is split into four sub-packages:

- `mpeg/ts` for parsing and processing MPEG-2 Transport Streams
- `mpeg/ps` for parsing and processing MPEG-2 Program Streams
- `mpeg/pes` for parsing and processing MPEG-2 Packetized Elementary Streams
- `mpeg/video` for decoding MPEG-2 Video

## Examples

### Give it a spin!

```
go get -d github.com/32bitkid/mpeg
```

### Decode a frame of video from a MPEG-2 TS and save it as a png

```go
package main

import "os"
import "image/png"

import "github.com/32bitkid/mpeg/ts"
import "github.com/32bitkid/mpeg/pes"
import "github.com/32bitkid/mpeg/video"

func main() {
  // Open the file
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
