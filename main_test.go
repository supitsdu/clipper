package main

import (
	"os"
	"testing"

	"github.com/atotto/clipboard"
)

func TestCopyToClipboard(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping clipboard test in CI/CD.")
	}

	content := "Hello, World!"
	err := copyToClipboard(content)
	if err != nil {
		t.Errorf("Failed to copy to clipboard: %v", err)
	}

	copied, err := clipboard.ReadAll()
	if err != nil {
		t.Errorf("Failed to read from clipboard: %v", err)
	}

	if copied != content {
		t.Errorf("Expected %s but got %s", content, copied)
	}
}

func TestReadFromStdin(t *testing.T) {
	content := "Stdin content"

	// Redirect stdin
	stdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.Write([]byte(content))
		w.Close()
	}()

	readContent, err := readFromStdin()
	if err != nil {
		t.Fatalf("Failed to read from stdin: %v", err)
	}

	os.Stdin = stdin

	if readContent != content {
		t.Errorf("Expected %s but got %s", content, readContent)
	}
}

func TestReadFromFile(t *testing.T) {
	content := "File content"
	file, err := os.CreateTemp("", "clipper")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	file.Close()

	readContent, err := readFromFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to read from file: %v", err)
	}

	if readContent != content {
		t.Errorf("Expected %s but got %s", content, readContent)
	}
}

func TestCopyFromCommandLineArgument(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping clipboard test in CI/CD.")
	}

	content := "Command line text"
	os.Args = []string{"clipper", "-c", content}

	main()

	copied, err := clipboard.ReadAll()
	if err != nil {
		t.Errorf("Failed to read from clipboard: %v", err)
	}

	if copied != content {
		t.Errorf("Expected %s but got %s", content, copied)
	}
}

func TestParseContentDirectText(t *testing.T) {
	directText := "Direct text input"
	contentStr, err := parseContent(&directText, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if contentStr != directText {
		t.Errorf("Expected contentStr to be '%s', got '%s'", directText, contentStr)
	}
}

func TestParseContentStdin(t *testing.T) {
	stdinInput := "Stdin input\n"
	r, w, _ := os.Pipe()
	originalStdin := os.Stdin
	defer func() {
		os.Stdin = originalStdin
		r.Close()
		w.Close()
	}()
	os.Stdin = r

	// Write the input to the pipe
	w.WriteString(stdinInput)
	w.Close()

	contentStr, err := parseContent(nil, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if contentStr != stdinInput {
		t.Errorf("Expected contentStr to be '%s', got '%s'", stdinInput, contentStr)
	}
}

func TestParseContentFromFile(t *testing.T) {
	contentFromFile := "Content from file"
	fakeFile := createTempFile(contentFromFile)
	defer os.Remove(fakeFile.Name())

	contentStr, err := parseContent(nil, []string{fakeFile.Name()})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if contentStr != contentFromFile+"\n" {
		t.Errorf("Expected contentStr to be '%s', got '%s'", contentFromFile, contentStr)
	}
}

// Helper function to create temporary file for testing
func createTempFile(content string) *os.File {
	file, _ := os.CreateTemp("", "testfile")
	file.WriteString(content)
	file.Close()
	return file
}
