package constants

import "errors"

var (
	ErrNotTSV   = errors.New("not a tsv file")
	ErrNotFound = errors.New("not found")
)
