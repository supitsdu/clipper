package reader

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/supitsdu/clipper/cli/options"
)

// ContentReader is responsible for reading and processing content from files or standard input.
type ContentReader struct {
	Config *options.Config // Configuration options for content reading and formatting.
}

// ReadAll reads content from multiple files concurrently or from standard input if no files are specified.
// It returns the aggregated content as a single string.
func (cr ContentReader) ReadAll() (string, error) {
	paths := cr.Config.FilePaths

	// If no file paths are specified, read from standard input.
	if len(paths) == 0 {
		return cr.IOReader(os.Stdin, "")
	}

	// Read content from all specified files concurrently.
	results, err := cr.ReadFilesAsync(paths)
	if err != nil {
		return "", err
	}

	// Join all results into a single string.
	return cr.JoinAll(results), nil
}

// ReadFilesAsync reads content from multiple files concurrently using goroutines.
// It returns a slice of strings containing the content of each file and an error if any file reading fails.
func (cr ContentReader) ReadFilesAsync(paths []string) ([]string, error) {
	errChan := make(chan error, len(paths)) // Channel to capture errors from goroutines
	results := make([]string, len(paths))   // Slice to store results

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, filepath := range paths {
		wg.Add(1)

		go func(i int, filepath string) { // Goroutine to read a single file.
			defer wg.Done()
			content, err := cr.ReadFile(filepath)

			if err != nil {
				// Send the error to the error channel and return early.
				errChan <- fmt.Errorf("error reading file '%s': %w", filepath, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()    // Ensure safe concurrent access to 'results'.
			results[i] = content // Store the content in the results slice.
		}(i, filepath)
	}

	go func() {
		wg.Wait() // Wait for all reading goroutines to complete.
		close(errChan)
	}()

	// Collect the first error encountered if any.
	var err error
	for e := range errChan {
		if err == nil {
			err = e
		}
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}

// JoinAll aggregates the content from the provided results and returns it as a single string.
// It combines the content of all readers into a single string with newline separators.
func (cr ContentReader) JoinAll(results []string) string {
	var sb strings.Builder
	for _, content := range results {
		if content != "" { // Ensure non-empty content is aggregated.
			sb.WriteString(content + "\n")
		}
	}
	return sb.String()
}

// ReadFile reads and formats content from a single file.
// It returns the formatted content as a string and an error if the file cannot be read or formatted.
func (cr ContentReader) ReadFile(filepath string) (string, error) {
	if err := cr.Readable(filepath); err != nil {
		return "", err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return cr.IOReader(file, filepath)
}

// Readable checks if a file path exists, is a regular file, and has read permissions.
// It returns an error if the file is not accessible or readable.
func (cr ContentReader) Readable(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// File doesn't exist or can't be accessed.
		return fmt.Errorf("file does not exist or can't be accessed")
	}

	// Check if it's a regular file (not a directory or other type).
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("path is not of a regular file, perhaps a directory or other type")
	}

	// Check if it's readable.
	if fileInfo.Mode().Perm()&0400 == 0 { // 0400 corresponds to read permission.
		return fmt.Errorf("you don't have access to read the file")
	}

	return nil
}

// IOReader reads content from an io.Reader (e.g., a file or stdin).
// It returns the content as a string and an error if reading fails.
func (cr ContentReader) IOReader(source io.Reader, filepath string) (string, error) {
	data, err := io.ReadAll(source)
	if err != nil {
		return "", err
	}

	return cr.CreateContent(filepath, data)
}

// CreateContent creates a formatted string from raw data and applies formatting based on the configuration.
func (cr ContentReader) CreateContent(filepath string, data []byte) (string, error) {
	// If no formatting is required, return the raw data as a string.
	if !cr.Config.ShouldFormat {
		return string(data), nil
	}

	return cr.Format(filepath, data)
}

// Format formats the content based on configuration options (HTML, Markdown, MimeType).
// It returns the formatted content as a string and an error if formatting fails.
func (cr ContentReader) Format(filepath string, data []byte) (string, error) {
	var sb strings.Builder

	mime := mimetype.Detect(data)
	mimeType := mime.String()
	content := string(data)

	if filepath == "" {
		filepath = "standard input"
	}

	// Append MIME type information if required.
	if cr.Config.MimeType {
		mimeType := fmt.Sprintf("%s (%s)\n", filepath, mimeType)

		if cr.Config.Html || cr.Config.Markdown {
			sb.WriteString("<!--\n" + mimeType + "-->\n")
		} else {
			sb.WriteString(mimeType)
		}
	}

	// Format the content based on the configuration.
	if cr.Config.Html {
		sb.WriteString(fmt.Sprintf("<code>\n%s\n</code>", content))
	} else if cr.Config.Markdown {
		sb.WriteString(fmt.Sprintf("```\n%s\n```", content))
	} else if cr.Config.LineNumbers {
		lines := strings.Split(content, "\n")

		for i, line := range lines {
			sb.WriteString(fmt.Sprintf("%d: %s\n", i+1, line))
		}
	} else {
		sb.WriteString(content)
	}

	return sb.String(), nil
}
