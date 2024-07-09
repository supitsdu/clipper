package reader_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supitsdu/clipper/cli/options"
	"github.com/supitsdu/clipper/cli/reader"
	"github.com/supitsdu/clipper/tests"
)

func TestReadAll_Files(t *testing.T) {
	content1 := "file1 content"
	content2 := "file2 content"

	file1, err := tests.CreateTempFile(t, content1)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	file2, err := tests.CreateTempFile(t, content2)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	config := &options.Config{
		FilePaths: []string{file1.Name(), file2.Name()},
	}
	contentReader := reader.ContentReader{Config: config}

	result, err := contentReader.ReadAll()
	assert.NoError(t, err)
	expectedResult := content1 + "\n" + content2 + "\n"
	assert.Equal(t, expectedResult, result)
}

func TestReadFile(t *testing.T) {
	content := "file content"
	file, err := tests.CreateTempFile(t, content)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	result, err := contentReader.ReadFile(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, content, result)
}

func TestReadable(t *testing.T) {
	file, err := tests.CreateTempFile(t, "content")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	err = contentReader.Readable(file.Name())
	assert.NoError(t, err)

	err = contentReader.Readable("nonexistentfile")
	assert.Error(t, err)
}

func TestIOReader(t *testing.T) {
	config := &options.Config{}
	contentReader := reader.ContentReader{Config: config}

	data := "some content"
	reader := strings.NewReader(data)

	result, err := contentReader.IOReader(reader, "")
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestCreateContent(t *testing.T) {
	t.Run("Without Content Formats", func(tt *testing.T) {
		config := &options.Config{ShouldFormat: false}
		contentReader := reader.ContentReader{Config: config}

		data := "raw content"
		result, err := contentReader.CreateContent("", []byte(data))
		assert.NoError(tt, err)
		assert.Equal(tt, data, result)
	})

	t.Run("Content with HTML5 format", func(tt *testing.T) {
		config := &options.Config{ShouldFormat: true, Html: true}
		contentReader := reader.ContentReader{Config: config}

		data := "raw content"
		result, err := contentReader.CreateContent("", []byte(data))
		assert.NoError(tt, err)
		assert.Equal(tt, "<code>\nraw content\n</code>", result)
	})

	t.Run("Content with Markdown format", func(tt *testing.T) {
		config := &options.Config{ShouldFormat: true, Markdown: true}
		contentReader := reader.ContentReader{Config: config}

		data := "raw content"
		result, err := contentReader.CreateContent("", []byte(data))
		assert.NoError(t, err)
		assert.Equal(t, "```\nraw content\n```", result)
	})
}

func TestJoinAll(t *testing.T) {
	t.Run("join contents with a new line", func(tt *testing.T) {
		config := &options.Config{}
		contentReader := reader.ContentReader{Config: config}

		results := []string{"content1", "content2", "content3"}
		expected := "content1\ncontent2\ncontent3\n"
		result := contentReader.JoinAll(results)

		assert.Equal(tt, expected, result)
	})

	t.Run("ignore empty contents", func(tt *testing.T) {
		config := &options.Config{}
		contentReader := reader.ContentReader{Config: config}

		results := []string{"", "content2", ""}
		expected := "content2\n"
		result := contentReader.JoinAll(results)

		assert.Equal(tt, expected, result)
	})
}