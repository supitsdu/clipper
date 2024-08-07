package console

import "errors"

var (
	ErrFileNotFound       = errors.New("file does not exist or cannot be accessed")
	ErrReadingDirectories = errors.New("reading from directories is not currently supported")
	ErrPermissionDenied   = errors.New("permission denied")
)
