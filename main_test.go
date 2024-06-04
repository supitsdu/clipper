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
