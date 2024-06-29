package clipper

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/supitsdu/clipper/cli/clipper"
	"github.com/xyproto/randomstring"
)

const mockTextContent = "Mocking Bird! Just A Sample Text."

func TestContentReaders(t *testing.T) {
	t.Run("FileContentReader", func(t *testing.T) {
		// Create a temporary file
		file, err := createTempFile(t, mockTextContent)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Create a FileContentReader
		reader := clipper.FileContentReader{FilePath: file.Name()}

		// Read the content
		readContent, err := reader.Read()
		if err != nil {
			t.Fatalf("Error reading file: %v", err)
		}

		// Check the content
		if readContent != mockTextContent {
			t.Errorf("Expected '%s', got '%s'", mockTextContent, readContent)
		}
	})

	t.Run("StdinContentReader", func(t *testing.T) {
		// Create a StdinContentReader
		reader := clipper.StdinContentReader{}

		// Replace stdin with a pipe
		_, w := replaceStdin(t)

		// Write some content to the pipe
		_, err := w.WriteString(mockTextContent)
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
		if readContent != mockTextContent {
			t.Errorf("Expected '%s', got '%s'", mockTextContent, readContent)
		}
	})
}

func TestClipboardWriter(t *testing.T) {
	t.Run("DefaultClipboardWriter", func(t *testing.T) {
		if testing.Short() == true {
			t.Skip("Skipping clipboard test in short mode. Helps avoid errors when on CI environments.")
		}

		// Create a DefaultClipboardWriter
		writer := clipper.DefaultClipboardWriter{}

		// Write some content to the clipboard
		err := writer.Write(mockTextContent)
		if err != nil {
			t.Errorf("Error writing to clipboard: %v", err)
		}

		// Check the clipboard content
		clipboardContent, err := clipboard.ReadAll()
		if err != nil {
			t.Errorf("Error reading from clipboard: %v", err)
		}

		// Check the content
		if clipboardContent != mockTextContent {
			t.Errorf("Expected '%s', got '%s'", mockTextContent, clipboardContent)
		}
	})
}

func TestReadContent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Single line content",
			input:       "Hello, World!",
			expected:    "Hello, World!",
			expectError: false,
		},
		{
			name:        "Multiple lines content",
			input:       "Hello, World!\nThis is a test.",
			expected:    "Hello, World!\nThis is a test.",
			expectError: false,
		},
		{
			name:        "Empty content",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "Content with EOF error",
			input:       "Hello, World!",
			expected:    "Hello, World!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			output, err := clipper.ReadContent(reader)
			if (err != nil) != tt.expectError {
				t.Fatalf("ReadContent() error = %v, expectError %v", err, tt.expectError)
			}
			if output != tt.expected {
				t.Errorf("ReadContent() = %v, want %v", output, tt.expected)
			}
		})
	}
}

func TestParseContent(t *testing.T) {
	t.Run("MultipleFiles", func(t *testing.T) {
		// Prepare test data with content from two files
		content1 := "Content from file 1"
		content2 := "Content from file 2"
		file1, err := createTempFile(t, content1)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		file2, err := createTempFile(t, content2)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Expected result after parsing multiple files
		expected := "Content from file 1\nContent from file 2\n"

		// Create FileContentReaders for the files
		reader1 := clipper.FileContentReader{FilePath: file1.Name()}
		reader2 := clipper.FileContentReader{FilePath: file2.Name()}

		// Execute the function under test
		actual, err := clipper.ParseContent(nil, reader1, reader2)
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
		emptyFile, err := createTempFile(t, "")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Expected result when parsing an empty file
		expected := "\n"

		// Create a FileContentReader for the empty file
		reader := clipper.FileContentReader{FilePath: emptyFile.Name()}

		// Execute the function under test with the empty file
		actual, err := clipper.ParseContent(nil, reader)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify the actual output matches the expected output
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})

	t.Run("InvalidNilInput", func(t *testing.T) {
		// Execute the function under test with no input arguments
		_, err := clipper.ParseContent(nil)
		if err == nil {
			t.Fatalf("Expected error, got %v", err)
		}
	})

	t.Run("InvalidFilePath", func(t *testing.T) {
		// Invalid file path that doesn't exist
		invalidPath := "/invalid/path/to/file.txt"

		// Create a FileContentReader with the invalid file path
		reader := clipper.FileContentReader{FilePath: invalidPath}

		// Execute the function under test with the invalid file path
		_, err := clipper.ParseContent(nil, reader)
		if err == nil {
			t.Fatalf("Expected an error, but got none")
		}
		// Further assertions can be added based on specific error expectations
		// For example, checking the error type or message.
	})

	t.Run("LargeFile", func(t *testing.T) {
		// Create a large file (e.g., 10MB)
		largeContent := randomstring.String(10 * 1024 * 1024) // 10MB of random data
		largeFile, err := createTempFile(t, largeContent)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		expected := largeContent + "\n"

		// Create a FileContentReader for the large file
		reader := clipper.FileContentReader{FilePath: largeFile.Name()}

		// Execute the function under test with the large file
		actual, err := clipper.ParseContent(nil, reader)
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
		directText := mockTextContent
		content, err := clipper.ParseContent(&directText)
		if err != nil {
			t.Errorf("Error parsing content: %v", err)
		}
		if content != directText {
			t.Errorf("Expected '%s', got '%s'", directText, content)
		}
	})

	t.Run("FileContentReader", func(t *testing.T) {
		// Test with a FileContentReader
		tmpFile, err := createTempFile(t, mockTextContent)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		reader := clipper.FileContentReader{FilePath: tmpFile.Name()}

		content, err := clipper.ParseContent(nil, reader)
		if err != nil {
			t.Errorf("Error parsing content: %v", err)
		}
		if content != (mockTextContent + "\n") {
			t.Errorf("Expected '%s', got '%s'", mockTextContent, content)
		}
	})
}

// replaceStdin replaces os.Stdin with a pipe for testing purposes.
func replaceStdin(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("replaceStdin: Failed to create pipe: %v", err)
	}

	// Save the original stdin
	originalStdin := os.Stdin

	// Replace stdin with the read end of the pipe
	os.Stdin = r

	// Cleanup function to restore the original stdin and close the pipe
	t.Cleanup(func() {
		os.Stdin = originalStdin
		r.Close()
		w.Close()
	})

	return r, w
}

// createTempFile creates a temporary file for testing purposes and writes the given content to it.
func createTempFile(t *testing.T, content string) (*os.File, error) {
	t.Helper()
	// Create a temporary file
	file, err := os.CreateTemp(t.TempDir(), "testfile")
	if err != nil {
		return nil, err
	}

	// Write the provided content to the temporary file
	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}
