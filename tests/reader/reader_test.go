package reader_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supitsdu/clipper/cli/options"
	"github.com/supitsdu/clipper/cli/reader"
	"github.com/supitsdu/clipper/tests"
)

func TestReadAll(t *testing.T) {
	t.Run("reads single file", func(t *testing.T) {
		file1, err := tests.CreateTempFile(t, tests.SampleText32)
		if err != nil {
			t.Fatalf("Failed to create temporary file: %s", err)
		}

		config := &options.Config{
			FilePaths: []string{file1.Name()},
		}

		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.ReadAll()
		assert.NoError(t, err)
		expectedResult := tests.SampleText32 + "\n"
		assert.Equal(t, expectedResult, result)
	})

	t.Run("reads from argument", func(t *testing.T) {
		config := &options.Config{
			Text: tests.SampleText32,
		}

		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.ReadAll()
		assert.NoError(t, err)
		expectedResult := tests.SampleText32 + "\n"
		assert.Equal(t, expectedResult, result)
	})

	t.Run("reads multiple files", func(t *testing.T) {
		fileContents := []string{"Japan", "Australia", "Germany"}
		expectedResults := strings.Join(fileContents, "\n") + "\n"
		var filePaths []string

		for index := 0; index < 3; index++ {
			file, err := tests.CreateTempFile(t, fileContents[index])
			if err != nil {
				t.Fatalf("Failed to create temporary file: %s", err)
			}

			filePaths = append(filePaths, file.Name())
		}

		config := &options.Config{FilePaths: filePaths}

		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.ReadAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedResults, result)
	})
}

func TestReadable(t *testing.T) {
	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	t.Run("file exists and is readable", func(t *testing.T) {
		file, err := tests.CreateTempFile(t, "content")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %s", err)
		}

		err = contentReader.Readable(file.Name())
		assert.NoError(t, err)
	})

	t.Run("file does not exist", func(t *testing.T) {
		err := contentReader.Readable("nonexistentfile")
		assert.Error(t, err)
		assert.Equal(t, "file does not exist or can't be accessed", err.Error())
	})

	t.Run("path is a directory", func(t *testing.T) {
		dir, err := tests.CreateTempDir(t, "testdir")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %s", err)
		}

		err = contentReader.Readable(dir)
		assert.Error(t, err)
		assert.Equal(t, "reading from directories is not currently supported", err.Error())
	})

	t.Run("file is not readable", func(t *testing.T) {
		file, err := tests.CreateTempFile(t, "content")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %s", err)
		}

		// Remove read permissions
		if err := os.Chmod(file.Name(), 0200); err != nil { // 0200 corresponds to write-only permission.
			t.Fatalf("Failed to change file permissions: %s", err)
		}

		err = contentReader.Readable(file.Name())
		assert.Error(t, err)
		assert.Equal(t, "permission denied", err.Error())

		// Restore read permissions for cleanup
		if err := os.Chmod(file.Name(), 0600); err != nil { // 0600 corresponds to read-write permissions.
			t.Fatalf("Failed to restore file permissions: %s", err)
		}
	})

	t.Run("file is readable after fixing permissions", func(t *testing.T) {
		file, err := tests.CreateTempFile(t, "content")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %s", err)
		}

		// Remove read permissions
		if err := os.Chmod(file.Name(), 0200); err != nil { // 0200 corresponds to write-only permission.
			t.Fatalf("Failed to change file permissions: %s", err)
		}

		err = contentReader.Readable(file.Name())
		assert.Error(t, err)
		assert.Equal(t, "permission denied", err.Error())

		// Restore read permissions
		if err := os.Chmod(file.Name(), 0600); err != nil { // 0600 corresponds to read-write permissions.
			t.Fatalf("Failed to restore file permissions: %s", err)
		}

		err = contentReader.Readable(file.Name())
		assert.NoError(t, err)
	})
}

func TestIOReader(t *testing.T) {
	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	t.Run("reads content from reader", func(t *testing.T) {
		reader := strings.NewReader(tests.SampleText32)
		result, err := contentReader.IOReader(reader, "")
		assert.NoError(t, err)
		assert.Equal(t, tests.SampleText32, result)
	})

	t.Run("returns empty string for empty reader", func(t *testing.T) {
		reader := strings.NewReader("")
		result, err := contentReader.IOReader(reader, "")
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("returns error for faulty reader", func(t *testing.T) {
		reader := io.NopCloser(&tests.FaultyReader{})
		result, err := contentReader.IOReader(reader, "")
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

}

func TestCreateContent(t *testing.T) {
	t.Run("Without Content Formats", func(t *testing.T) {
		config := &options.Config{ShouldFormat: false}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.CreateContent("", []byte(tests.SampleText32))
		assert.NoError(t, err)
		assert.Equal(t, tests.SampleText32, result)
	})

	t.Run("Content with HTML5 format", func(t *testing.T) {
		config := &options.Config{ShouldFormat: true, HTML: true}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.CreateContent("", []byte(tests.SampleText32))
		assert.NoError(t, err)
		assert.Equal(t, "<code>\n"+tests.SampleText32+"\n</code>", result)
	})

	t.Run("Content with Markdown format", func(t *testing.T) {
		config := &options.Config{ShouldFormat: true, Markdown: true}
		contentReader := reader.ContentReader{Config: config}

		result, err := contentReader.CreateContent("", []byte(tests.SampleText32))
		assert.NoError(t, err)
		assert.Equal(t, "```\n"+tests.SampleText32+"\n```", result)
	})
}

func TestJoinAll(t *testing.T) {
	t.Run("join contents with a new line", func(t *testing.T) {
		config := &options.Config{}
		contentReader := reader.ContentReader{Config: config}

		results := []string{"content1", "content2", "content3"}
		expected := "content1\ncontent2\ncontent3\n"
		result := contentReader.JoinAll(results)

		assert.Equal(t, expected, result)
	})

	t.Run("ignore empty contents", func(t *testing.T) {
		config := &options.Config{}
		contentReader := reader.ContentReader{Config: config}

		results := []string{"", "content2", ""}
		expected := "content2\n"
		result := contentReader.JoinAll(results)

		assert.Equal(t, expected, result)
	})
}
