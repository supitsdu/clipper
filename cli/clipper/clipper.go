package clipper

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/supitsdu/clipper/cli/options"
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
	file, err := os.Open(f.FilePath)
	if err != nil {
		return "", fmt.Errorf("error opening file '%s': %w", f.FilePath, err)
	}
	defer file.Close()

	return ReadContent(bufio.NewReader(file))
}

// StdinContentReader reads content from the standard input (stdin).
type StdinContentReader struct{}

// Read reads the content from stdin.
func (s StdinContentReader) Read() (string, error) {
	return ReadContent(bufio.NewReader(os.Stdin))
}

// DefaultClipboardWriter writes content to the clipboard using the default clipboard implementation.
type DefaultClipboardWriter struct{}

// Write writes the given content to the clipboard.
func (c DefaultClipboardWriter) Write(content string) error {
	return clipboard.WriteAll(content)
}

// readContent reads content from the provided bufio.Reader and returns it as a string.
func ReadContent(reader *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				sb.WriteString(line)
				break
			}
			return "", fmt.Errorf("error reading content: %w", err)
		}
		sb.WriteString(line)
	}
	return sb.String(), nil
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

func GetReaders(targets []string) []ContentReader {
	if len(targets) == 0 {
		// If no file paths are provided, use StdinContentReader to read from stdin.
		return []ContentReader{StdinContentReader{}}
	}

	// If file paths are provided as arguments, create FileContentReader instances for each.
	var readers []ContentReader
	for _, filePath := range targets {
		readers = append(readers, FileContentReader{FilePath: filePath})
	}
	return readers
}

// Run executes the clipper tool logic based on the provided configuration.
func Run(config *options.Config, writer ClipboardWriter) (string, error) {
	if *config.ShowVersion {
		return options.Version, nil
	}

	readers := GetReaders(config.Args)

	// Aggregate the content from the provided sources.
	content, err := ParseContent(config.DirectText, readers...)
	if err != nil {
		return "", fmt.Errorf("parsing content: %w", err)
	}

	// Write the parsed content to the provided clipboard.
	if err = writer.Write(content); err != nil {
		return "", fmt.Errorf("copying content to clipboard: %w", err)
	}

	return "updated clipboard successfully. Ready to paste!", nil
}
