package ps

import "errors"

var (
	ErrPackStartCodeNotFound  = errors.New("pack start code not found")
	ErrProgramEndCodeNotFound = errors.New("program end code not found")
	ErrMarkerNotFound         = errors.New("marker not found")
)
