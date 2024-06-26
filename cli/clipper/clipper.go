package clipper

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/supitsdu/clipper/cli/options"
	"github.com/atotto/clipboard"
)

// ContentReader defines an interface for reading content from various sources.
type ContentReader interface {
	Read() (string, error)
}

// ClipboardWriter defines an interface for writing content to the clipboard.
type ClipboardWriter interface {
	Write(content string) error
}

// FileContentReader reads content from a specified file path.
type FileContentReader struct {
	FilePath string
}

// Read reads the content from the file specified in FileContentReader.
func (f FileContentReader) Read() (string, error) {
	content, err := os.ReadFile(f.FilePath)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %v", f.FilePath, err)
	}
	return string(content), nil
}

// StdinContentReader reads content from the standard input (stdin).
type StdinContentReader struct{}

// Read reads the content from stdin.
func (s StdinContentReader) Read() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("error reading from stdin: %v", err)
	}
	return string(input), nil
}

// DefaultClipboardWriter writes content to the clipboard using the default clipboard implementation.
type DefaultClipboardWriter struct{}

// Write writes the given content to the clipboard.
func (c DefaultClipboardWriter) Write(content string) error {
	return clipboard.WriteAll(content)
}

// ParseContent aggregates content from the provided readers, or returns the direct text if provided.
func ParseContent(directText *string, readers ...ContentReader) (string, error) {
	if directText != nil && *directText != "" {
		return *directText, nil
	}

	if len(readers) == 0 {
		return "", fmt.Errorf("no content readers provided")
	}

	var sb strings.Builder
	for _, reader := range readers {
		content, err := reader.Read()
		if err != nil {
			return "", err
		}
		sb.WriteString(content + "\n")
	}

	return sb.String(), nil
}

// Run executes the clipper tool logic based on the provided configuration.
func Run(config *options.Config) {
	// Display the version if the flag is set.
	if *config.ShowVersion {
		fmt.Printf("Clipper %s\n", options.Version)
		return
	}

	var readers []ContentReader
	if len(config.Args) > 0 {
		// If file paths are provided as arguments, create FileContentReader instances for each.
		for _, filePath := range config.Args {
			readers = append(readers, FileContentReader{FilePath: filePath})
		}
	} else {
		// If no file paths are provided, use StdinContentReader to read from stdin.
		readers = append(readers, StdinContentReader{})
	}

	// Aggregate the content from the provided sources.
	content, err := ParseContent(config.DirectText, readers...)
	if err != nil {
		fmt.Printf("Error parsing content: %v\n", err)
		os.Exit(1)
	}

	// Write the parsed content to the clipboard.
	writer := DefaultClipboardWriter{}
	if err = writer.Write(content); err != nil {
		fmt.Printf("Error copying content to clipboard: %v\n", err)
		os.Exit(1)
	}

	// Print success message.
	fmt.Println("Clipboard updated successfully. Ready to paste!")
}
