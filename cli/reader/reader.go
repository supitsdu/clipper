package reader

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// ContentReader defines an interface for reading content from various sources.
type ContentReader interface {
	Read() (string, error)
}

// FileContentReader reads content from a specified file path.
type FileContentReader struct {
	FilePath string
}

// Read reads the content from the file specified in FileContentReader.
// It reads the entire file content into memory, which is suitable for smaller files.
func (f FileContentReader) Read() (string, error) {
	content, err := os.ReadFile(f.FilePath)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %w", f.FilePath, err)
	}
	return string(content), nil
}

// StdinContentReader reads content from the standard input (stdin).
type StdinContentReader struct{}

// Read reads the content from stdin.
// It reads all the data from stdin until EOF, which is useful for piping input.
func (s StdinContentReader) Read() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}
	return string(input), nil
}

// ReadContentConcurrently reads content from multiple readers concurrently and returns the results.
// This function utilizes goroutines to perform concurrent reads, improving performance for multiple files.
func ReadContentConcurrently(readers []ContentReader) ([]string, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(readers)) // Channel to capture errors
	results := make([]string, len(readers))   // Slice to store results

	for i, reader := range readers {
		wg.Add(1)
		go func(i int, reader ContentReader) {
			defer wg.Done()
			content, err := reader.Read()
			if err != nil {
				errChan <- err // Send error to channel
				return
			}
			mu.Lock()
			results[i] = content // Safely write to results slice
			mu.Unlock()
		}(i, reader)
	}

	wg.Wait()
	close(errChan) // Close error channel after all reads are done

	if len(errChan) > 0 {
		return nil, <-errChan // Return the first error encountered
	}

	return results, nil
}

// AggregateContent aggregates the content from the provided results and returns it as a single string.
// It combines the content of all readers into a single string with newline separators.
func AggregateContent(results []string) string {
	var sb strings.Builder
	for _, content := range results {
		sb.WriteString(content + "\n")
	}
	return sb.String()
}

// ParseContent aggregates content from the provided readers, or returns the direct text if provided.
// This function first checks for direct text input, then reads from the provided readers concurrently.
func ParseContent(directText *string, readers ...ContentReader) (string, error) {
	if directText != nil && *directText != "" {
		return *directText, nil // Return direct text if provided
	}

	if len(readers) == 0 {
		return "", fmt.Errorf("no content readers provided")
	}

	results, err := ReadContentConcurrently(readers) // Read content concurrently
	if err != nil {
		return "", err
	}

	return AggregateContent(results), nil // Aggregate and return the content
}

// GetReaders constructs the appropriate ContentReaders based on the provided file paths or lack thereof.
// If no targets are provided, it defaults to using StdinContentReader.
func GetReaders(targets []string) []ContentReader {
	if len(targets) == 0 {
		// If no file paths are provided, use StdinContentReader to read from stdin.
		return []ContentReader{StdinContentReader{}}
	}

	// Create FileContentReader instances for each provided file path.
	var readers []ContentReader
	for _, filePath := range targets {
		readers = append(readers, FileContentReader{FilePath: filePath})
	}
	return readers
}
