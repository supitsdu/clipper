package main

import (
	"os"
	"testing"
	"github.com/xyproto/randomstring"
)

// TestParseContent_MultipleFiles verifies parseContent correctly concatenates contents of multiple files.
func TestParseContent_MultipleFiles(t *testing.T) {
	// Prepare test data with content from two files
	content1 := "Content from file 1"
	content2 := "Content from file 2"
	file1 := createTempFile(t, content1)
	file2 := createTempFile(t, content2)

	// Expected result after parsing multiple files
	expected := "Content from file 1\nContent from file 2\n"

	// Define directText as nil since we are testing file paths only
	var directText *string

	// Execute the function under test
	actual, err := parseContent(directText, []string{file1.Name(), file2.Name()})
	if err != nil {
		t.Fatalf("TestParseContent_MultipleFiles: Expected no error, got %v", err)
	}

	// Verify the actual output matches the expected output
	if actual != expected {
		t.Errorf("TestParseContent_MultipleFiles: Expected %s but got %s", expected, actual)
	}
}

// TestParseContent_EmptyFiles verifies parseContent handles empty files gracefully.
func TestParseContent_EmptyFiles(t *testing.T) {
	// Create an empty temporary file
	emptyFile := createTempFile(t, "")

	// Expected result when parsing an empty file
	expected := "\n"

	// Define directText as nil since we are testing file paths only
	var directText *string

	// Execute the function under test with the empty file
	actual, err := parseContent(directText, []string{emptyFile.Name()})
	if err != nil {
		t.Fatalf("TestParseContent_EmptyFiles: Expected no error, got %v", err)
	}

	// Verify the actual output matches the expected output
	if actual != expected {
		t.Errorf("TestParseContent_EmptyFiles: Expected %s but got %s", expected, actual)
	}
}

// TestParseContent_EmptyInput ensures parseContent handles empty input gracefully.
func TestParseContent_EmptyInput(t *testing.T) {
	// Expected result for empty input
	expected := ""

	// Define directText as nil since we are testing empty inputs only
	var directText *string

	// Execute the function under test with no input arguments
	actual, err := parseContent(directText, []string{})
	if err != nil {
		t.Fatalf("TestParseContent_EmptyInput: Expected no error, got %v", err)
	}

	// Verify the actual output matches the expected output
	if actual != expected {
		t.Errorf("TestParseContent_EmptyInput: Expected %s but got %s", expected, actual)
	}
}

// TestParseContent_InvalidFilePath checks how parseContent behaves with an invalid file path.
func TestParseContent_InvalidFilePath(t *testing.T) {
	// Invalid file path that doesn't exist
	invalidPath := "/invalid/path/to/file.txt"

	// Define directText as nil since we are testing file paths only
	var directText *string

	// Execute the function under test with the invalid file path
	_, err := parseContent(directText, []string{invalidPath})
	if err == nil {
		t.Fatalf("TestParseContent_InvalidFilePath: Expected an error, but got none")
	}
	// Further assertions can be added based on specific error expectations
	// For example, checking the error type or message.
}

// TestParseContent_LargeFile verifies parseContent handles large files correctly
func TestParseContent_LargeFile(t *testing.T) {
	// Create a large file (e.g., 10MB)
	largeContent := randomstring.String(10 * 1024 * 1024) // 10MB of random data
	largeFile := createTempFile(t, largeContent)
	expected := largeContent + "\n"

	// Define directText as nil since we are testing file paths only
	var directText *string

	// Execute the function under test with the large file
	actual, err := parseContent(directText, []string{largeFile.Name()})
	if err != nil {
		t.Fatalf("TestParseContent_LargeFile: Expected no error, got %v", err)
	}

	// Verify the actual output matches the expected output
	if actual != expected {
		t.Errorf("TestParseContent_LargeFile: Output differs")
	}
}

// TestParseContent_DirectText verifies parseContent handles direct text input correctly.
func TestParseContent_DirectText(t *testing.T) {
	// Direct text input to parse
	content := "Direct text input"
	expected := "Direct text input"

	// Execute the function under test with direct text input
	actual, err := parseContent(stringPtr(content), []string{})
	if err != nil {
		t.Fatalf("TestParseContent_DirectText: Expected no error, got %v", err)
	}

	// Verify the actual output matches the expected output
	if actual != expected {
		t.Errorf("TestParseContent_DirectText: Expected %s but got %s", expected, actual)
	}
}

// TestParseContent_Stdin verifies parseContent handles stdin input correctly.
func TestParseContent_Stdin(t *testing.T) {
    // Stdin input content
    content := "Stdin input\n"
    expected := "Stdin input\n"
	// Define directText as nil since we are testing stdin only
	var directText *string

    // Replace stdin with a pipe for testing
    _, w := replaceStdin(t)

    // Write content to the pipe and close it
    go func() {
        _, err := w.Write([]byte(content))
        if err != nil {
            t.Fatalf("Failed to write to stdin pipe: %v", err)
        }
        w.Close()
    }()

    // Execute the function under test with stdin input
    actual, err := parseContent(directText, []string{})
    if err != nil {
        t.Fatalf("TestParseContent_Stdin: Expected no error, got %v", err)
    }

    // Verify the actual output matches the expected output
    if actual != expected {
        t.Errorf("TestParseContent_Stdin: Expected %s but got %s", expected, actual)
    }
}

// TestParseContent_File verifies parseContent handles file input correctly.
func TestParseContent_File(t *testing.T) {
	// Content to write to the temporary file
	content := "Content from file"
	file := createTempFile(t, content)

	// Expected content after parsing the file
	expected := "Content from file\n"

	// Execute the function under test with the temporary file as input
	actual, err := parseContent(nil, []string{file.Name()})
	if err != nil {
		t.Fatalf("TestParseContent_File: Expected no error, got %v", err)
	}

	// Verify the actual output matches the expected output
	if actual != expected {
		t.Errorf("TestParseContent_File: Expected %s but got %s", expected, actual)
	}
}

// stringPtr creates a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// createTempFile creates a temporary file for testing purposes and writes the given content to it.
func createTempFile(t *testing.T, content string) *os.File {
t.Helper ()
	// Create a temporary file
	file, err := os.CreateTemp(t.TempDir(), "testfile")
	if err != nil {
		t.Fatalf("createTempFile: Failed to create temp file: %v", err)
	}

	// Write the provided content to the temporary file
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("createTempFile: Failed to write to temp file: %v", err)
	}

	return file
}

// replaceStdin replaces os.Stdin with a pipe for testing purposes.
func replaceStdin(t *testing.T) (*os.File, *os.File) {
t.Helper ()
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
