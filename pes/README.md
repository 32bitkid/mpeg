# mpeg/pes

Package pes implements an MPEG-2 packetized elementary stream parser as defined in ISO/IEC 13818-1.

[![GoDoc](https://godoc.org/github.com/32bitkid/mpeg/pes?status.svg)](https://godoc.org/github.com/32bitkid/mpeg/pes)

This package is experimental and is not intended for
use in production environments.

## Roadmap
- [x] Basic packet parser
- [ ] Full PES Extension support
- [ ] Packet streamer interface (via `<-chan Packet`)
- [x] PayloadReader (implementing `io.Reader`)
