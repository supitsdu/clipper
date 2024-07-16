package reader_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/supitsdu/clipper/cli/options"
	"github.com/supitsdu/clipper/cli/reader"
	"github.com/supitsdu/clipper/tests"
)

const (
	ErrFileNotFound       = "file does not exist or can't be accessed"
	ErrReadingDirectories = "reading from directories is not currently supported"
	ErrPermissionDenied   = "permission denied"
)

func TestReadAll(t *testing.T) {
	t.Run("reads single file", func(t *testing.T) {
		file := tests.CreateTempFile(t, tests.SampleText32)

		config := &options.Config{
			FilePaths: []string{file.Name()},
		}

		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.ReadAll()
		require.NoError(t, err)
		expectedResult := tests.SampleText32 + "\n"
		require.Equal(t, expectedResult, result)
	})

	t.Run("reads from argument", func(t *testing.T) {
		config := &options.Config{
			Text: tests.SampleText32,
		}

		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.ReadAll()
		require.NoError(t, err)
		expectedResult := tests.SampleText32 + "\n"
		require.Equal(t, expectedResult, result)
	})

	t.Run("reads multiple files", func(t *testing.T) {
		fileContents := []string{"Japan", "Australia", "Germany"}
		expectedResults := strings.Join(fileContents, "\n") + "\n"
		var filePaths []string

		for index := 0; index < 3; index++ {
			file := tests.CreateTempFile(t, fileContents[index])

			filePaths = append(filePaths, file.Name())
		}

		config := &options.Config{FilePaths: filePaths}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.ReadAll()

		require.NoError(t, err)
		require.Equal(t, expectedResults, result)
	})
}

//nolint:funlen // It cannot be shorter
func TestReadable(t *testing.T) {
	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	t.Run("file exists and is readable", func(t *testing.T) {
		file := tests.CreateTempFile(t, "content")
		err := contentReader.Readable(file.Name())
		require.NoError(t, err)
	})

	t.Run("file does not exist", func(t *testing.T) {
		err := contentReader.Readable("nonexistentfile")
		require.Error(t, err)
		assert.EqualError(t, err, ErrFileNotFound)
	})

	t.Run("path is a directory", func(t *testing.T) {
		// Create a temporary directory
		dir, err := os.MkdirTemp(t.TempDir(), "testdir")
		require.NoError(t, err)

		err = contentReader.Readable(dir)
		require.Error(t, err)
		assert.EqualError(t, err, ErrReadingDirectories)
	})

	t.Run("file is not readable", func(t *testing.T) {
		file := tests.CreateTempFile(t, "content")

		// Remove read permissions
		if err := os.Chmod(file.Name(), 0o200); err != nil { // 0200 corresponds to write-only permission.
			t.Fatalf("Failed to change file permissions: %s", err)
		}

		err := contentReader.Readable(file.Name())
		require.Error(t, err)
		assert.EqualError(t, err, ErrPermissionDenied)

		// Restore read permissions for cleanup
		if err := os.Chmod(file.Name(), 0o600); err != nil { // 0600 corresponds to read-write permissions.
			t.Fatalf("Failed to restore file permissions: %s", err)
		}
	})

	t.Run("file is readable after fixing permissions", func(t *testing.T) {
		file := tests.CreateTempFile(t, "content")

		// Remove read permissions
		if err := os.Chmod(file.Name(), 0o200); err != nil { // 0200 corresponds to write-only permission.
			t.Fatalf("Failed to change file permissions: %s", err)
		}

		err := contentReader.Readable(file.Name())
		require.Error(t, err)
		assert.EqualError(t, err, ErrPermissionDenied)

		// Restore read permissions
		if err := os.Chmod(file.Name(), 0o600); err != nil { // 0600 corresponds to read-write permissions.
			t.Fatalf("Failed to restore file permissions: %s", err)
		}

		err = contentReader.Readable(file.Name())
		require.NoError(t, err)
	})
}

func TestIOReader(t *testing.T) {
	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	t.Run("reads content from reader", func(t *testing.T) {
		reader := strings.NewReader(tests.SampleText32)
		result, err := contentReader.IOReader(reader, "")
		require.NoError(t, err)
		require.Equal(t, tests.SampleText32, result)
	})

	t.Run("returns empty string for empty reader", func(t *testing.T) {
		reader := strings.NewReader("")
		result, err := contentReader.IOReader(reader, "")
		require.NoError(t, err)
		require.Equal(t, "", result)
	})

	t.Run("returns error for faulty reader", func(t *testing.T) {
		reader := io.NopCloser(&tests.FaultyReader{})
		result, err := contentReader.IOReader(reader, "")
		require.Error(t, err)
		require.Equal(t, "", result)
	})
}

func TestCreateContent(t *testing.T) {
	t.Run("Without Content Formats", func(t *testing.T) {
		config := &options.Config{ShouldFormat: false}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.CreateContent("", []byte(tests.SampleText32))
		require.NoError(t, err)
		require.Equal(t, tests.SampleText32, result)
	})

	t.Run("Content with HTML5 format", func(t *testing.T) {
		config := &options.Config{ShouldFormat: true, HTML: true}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.CreateContent("", []byte(tests.SampleText32))
		require.NoError(t, err)
		require.Equal(t, "<code>\n"+tests.SampleText32+"\n</code>", result)
	})

	t.Run("Content with Markdown format", func(t *testing.T) {
		config := &options.Config{ShouldFormat: true, Markdown: true}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.CreateContent("", []byte(tests.SampleText32))
		require.NoError(t, err)
		require.Equal(t, "```\n"+tests.SampleText32+"\n```", result)
	})
}

func TestJoinAll(t *testing.T) {
	t.Run("join contents with a new line", func(t *testing.T) {
		config := &options.Config{}
		contentReader := reader.ContentReader{Config: config}

		results := []string{"content1", "content2", "content3"}
		expected := "content1\ncontent2\ncontent3\n"
		result := contentReader.JoinAll(results)

		require.Equal(t, expected, result)
	})

	t.Run("ignore empty contents", func(t *testing.T) {
		config := &options.Config{}
		contentReader := reader.ContentReader{Config: config}

		results := []string{"", "content2", ""}
		expected := "content2\n"
		result := contentReader.JoinAll(results)

		require.Equal(t, expected, result)
	})
}
