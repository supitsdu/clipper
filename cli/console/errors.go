package console

import "fmt"

var (
	ErrFileNotFound       = fmt.Errorf("file does not exist or can't be accessed")
	ErrReadingDirectories = fmt.Errorf("reading from directories is not currently supported")
	ErrPermissionDenied   = fmt.Errorf("permission denied")
)
