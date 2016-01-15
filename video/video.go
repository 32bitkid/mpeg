// Package video implements the basic structures for simple MPEG-2 video decoding as defined in ISO/IEC 13818-1.
//
// This package is experimental and is not intended for
// use in production environments.
//
// Presently, this library supports decoding a subset of the
// entire MPEG-2 decoding specification: namely frame based pictures,
// subsampled by 4:2:0 can be decoded. However, this package is an active work
// in progress and slowly inching toward broader support the spec.
package video
