package reader_test

import (
	"testing"

	"github.com/supitsdu/clipper/cli/reader"
	"github.com/supitsdu/clipper/tests"
	"github.com/xyproto/randomstring"
)

func TestContentReaders(t *testing.T) {
	t.Run("FileContentReader", func(t *testing.T) {
		// Create a temporary file
		file, err := tests.CreateTempFile(t, tests.SampleText_32)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Create a FileContentReader
		reader := reader.FileContentReader{FilePath: file.Name()}

		// Read the content
		readContent, err := reader.Read()
		if err != nil {
			t.Fatalf("Error reading file: %v", err)
		}

		// Check the content
		if readContent != tests.SampleText_32 {
			t.Errorf("Expected '%s', got '%s'", tests.SampleText_32, readContent)
		}
	})

	t.Run("StdinContentReader", func(t *testing.T) {
		// Create a StdinContentReader
		reader := reader.StdinContentReader{}

		// Replace stdin with a pipe
		_, w := tests.ReplaceStdin(t)

		// Write some content to the pipe
		_, err := w.WriteString(tests.SampleText_32)
		if err != nil {
			t.Fatalf("Error writing to pipe: %v", err)
		}

		// Close the write end of the pipe
		err = w.Close()
		if err != nil {
			t.Fatalf("Error closing pipe: %v", err)
		}

		// Read the content
		readContent, err := reader.Read()
		if err != nil {
			t.Fatalf("Error reading from stdin: %v", err)
		}

		// Check the content
		if readContent != tests.SampleText_32 {
			t.Errorf("Expected '%s', got '%s'", tests.SampleText_32, readContent)
		}
	})
}

func TestParseContent(t *testing.T) {
	t.Run("MultipleFiles", func(t *testing.T) {
		// Prepare test data with content from two files
		content1 := "Content from file 1"
		content2 := "Content from file 2"
		file1, err := tests.CreateTempFile(t, content1)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		file2, err := tests.CreateTempFile(t, content2)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Expected result after parsing multiple files
		expected := "Content from file 1\nContent from file 2\n"

		// Create FileContentReaders for the files
		reader1 := reader.FileContentReader{FilePath: file1.Name()}
		reader2 := reader.FileContentReader{FilePath: file2.Name()}

		// Execute the function under test
		actual, err := reader.ParseContent("", reader1, reader2)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the actual output matches the expected output
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})

	t.Run("EmptyFiles", func(t *testing.T) {
		// Create an empty temporary file
		emptyFile, err := tests.CreateTempFile(t, "")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Expected result when parsing an empty file
		emptyString := ""
		newlineString := "\n"

		// Create a FileContentReader for the empty file
		testReader := reader.FileContentReader{FilePath: emptyFile.Name()}

		// Execute the function under test with the empty file
		actual, err := reader.ParseContent("", testReader)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if actual != emptyString {
			t.Error("When given an empty file returned something else.")
		}

		if actual == newlineString {
			t.Error("When given empty file returned one with newline '\\n' character. See: https://github.com/supitsdu/clipper/issues/31#issue-2379691170")
		}
	})

	t.Run("InvalidNilInput", func(t *testing.T) {
		// Execute the function under test with no input arguments
		_, err := reader.ParseContent("")
		if err == nil {
			t.Fatalf("Expected error, got %v", err)
		}
	})

	t.Run("InvalidFilePath", func(t *testing.T) {
		// Invalid file path that doesn't exist
		invalidPath := "/invalid/path/to/file.txt"

		// Create a FileContentReader with the invalid file path
		testReader := reader.FileContentReader{FilePath: invalidPath}

		// Execute the function under test with the invalid file path
		_, err := reader.ParseContent("", testReader)
		if err == nil {
			t.Fatalf("Expected an error, but got none")
		}
		// Further assertions can be added based on specific error expectations
		// For example, checking the error type or message.
	})

	t.Run("LargeFile", func(t *testing.T) {
		// Create a large file (e.g., 10MB)
		largeContent := randomstring.String(10 * 1024 * 1024) // 10MB of random data
		largeFile, err := tests.CreateTempFile(t, largeContent)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		expected := largeContent + "\n"

		// Create a FileContentReader for the large file
		testReader := reader.FileContentReader{FilePath: largeFile.Name()}

		// Execute the function under test with the large file
		actual, err := reader.ParseContent("", testReader)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the actual output matches the expected output
		if actual != expected {
			t.Errorf("Output differs")
		}
	})

	t.Run("DirectText", func(t *testing.T) {
		// Test with direct text
		directText := tests.SampleText_32
		content, err := reader.ParseContent(directText)
		if err != nil {
			t.Errorf("Error parsing content: %v", err)
		}
		if content != directText {
			t.Errorf("Expected '%s', got '%s'", directText, content)
		}
	})

	t.Run("FileContentReader", func(t *testing.T) {
		// Test with a FileContentReader
		tmpFile, err := tests.CreateTempFile(t, tests.SampleText_32)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		testReader := reader.FileContentReader{FilePath: tmpFile.Name()}

		content, err := reader.ParseContent("", testReader)
		if err != nil {
			t.Errorf("Error parsing content: %v", err)
		}
		if content != (tests.SampleText_32 + "\n") {
			t.Errorf("Expected '%s', got '%s'", tests.SampleText_32, content)
		}
	})
}
