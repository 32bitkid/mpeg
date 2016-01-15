# mpeg/ts

Package ts implements an MPEG-2 transport stream parser as defined in ISO/IEC 13818-1.

[![GoDoc](https://godoc.org/github.com/32bitkid/mpeg/ts?status.svg)](https://godoc.org/github.com/32bitkid/mpeg/ts)

## Roadmap

- [x] Basic packet parser
- [ ] Adapation feild support
- [x] Stream de-multiplexer (via `<-chan Packet`)
- [x] PayloadReader (implementing `io.Reader`)
- [x] PayloadUnitReader (implementing `io.Reader`)
- [ ] Decoding support for other TS Packet types
    - [ ] Program Association Table support
    - [ ] Conditional Access Table support
    - [ ] Program Map Table support

