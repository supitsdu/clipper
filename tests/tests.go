package tests

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const SampleText32 = "Mocking Bird! Just A Sample Text."

// replaceStdin replaces os.Stdin with a pipe for testing purposes.
func ReplaceStdin(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	// Create a pipe
	readEnd, writeEnd, err := os.Pipe()
	if err != nil {
		t.Fatalf("replaceStdin: Failed to create pipe: %v", err)
	}

	// Save the original stdin
	originalStdin := os.Stdin

	// Replace stdin with the read end of the pipe
	os.Stdin = readEnd

	// Cleanup function to restore the original stdin and close the pipe
	t.Cleanup(func() {
		os.Stdin = originalStdin
		readEnd.Close()
		writeEnd.Close()
	})

	return readEnd, writeEnd
}

// createTempFile creates a temporary file for testing purposes and writes the given content to it.
func CreateTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	// Create a temporary file
	file, err := os.CreateTemp(t.TempDir(), "testfile")
	require.NoError(t, err)

	// Write the provided content to the temporary file
	_, err = file.WriteString(content)
	require.NoError(t, err)

	return file
}

// IsCIEnvironment checks if the code is running in a CI environment.
func IsCIEnvironment() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// faultyReader simulates a reader that always returns an error.
type FaultyReader struct{}

func (*FaultyReader) Read(p []byte) (n int, err error) {
	return len(p), io.ErrUnexpectedEOF
}
