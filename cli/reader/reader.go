package reader

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/supitsdu/clipper/cli/console"
	"github.com/supitsdu/clipper/cli/options"
)

// ContentReader is responsible for reading and processing content from files or standard input.
type ContentReader struct {
	Config *options.Config // Configuration options for content reading and formatting.
}

// ReadAll reads content from multiple sources: command-line argument, standard input (stdin),
// and specified files. It aggregates the content into a single string.
// It returns the aggregated content as a string and any error encountered during the process.
func (cr ContentReader) ReadAll() (string, error) {
	var results []string

	// Add text from command-line argument, if provided
	if len(cr.Config.Text) > 0 {
		results = append(results, cr.Config.Text)
	}

	// Read from stdin if data is available
	if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		stdinContents, err := cr.IOReader(os.Stdin, "")
		if err != nil {
			return "", err
		}
		if len(stdinContents) > 0 {
			results = append(results, stdinContents)
		}
	}

	// Read from specified files asynchronously
	fileContents, err := cr.ReadFilesAsync(cr.Config.FilePaths)
	if err != nil {
		return "", err
	}

	results = append(results, fileContents...)

	// Join all results into a single string.
	return cr.JoinAll(results), nil
}

// ReadFilesAsync reads content from multiple files concurrently using goroutines.
// It returns a slice of strings containing the content of each file and an error if any file reading fails.
func (cr ContentReader) ReadFilesAsync(paths []string) ([]string, error) {
	errChan := make(chan error, len(paths)) // Channel to capture errors from goroutines
	results := make([]string, len(paths))   // Slice to store results

	var waitGroup sync.WaitGroup
	var mutext sync.Mutex

	for index, filepath := range paths {
		waitGroup.Add(1)

		go func(index int, filepath string) { // Goroutine to read a single file.
			defer waitGroup.Done()
			content, err := cr.ReadFile(filepath)

			if err != nil {
				// Send the error to the error channel and return early.
				errChan <- fmt.Errorf("reading file '%s': %w", filepath, err)
				return
			}

			mutext.Lock()
			defer mutext.Unlock()    // Ensure safe concurrent access to 'results'.
			results[index] = content // Store the content in the results slice.
		}(index, filepath)
	}

	go func() {
		waitGroup.Wait() // Wait for all reading goroutines to complete.
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
	var strBuilder strings.Builder
	for _, content := range results {
		if content != "" { // Ensure non-empty content is aggregated.
			strBuilder.WriteString(content + "\n")
		}
	}
	return strBuilder.String()
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
		return console.ErrFileNotFound
	}

	// Check if it's a regular file (not a directory or other type).
	if !fileInfo.Mode().IsRegular() {
		return console.ErrReadingDirectories
	}

	// Check if it's readable.
	if fileInfo.Mode().Perm()&0400 == 0 { // 0400 corresponds to read permission.
		return console.ErrPermissionDenied
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

	// Check if stdin is empty
	if len(data) == 0 {
		return "", nil
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
	var strBuilder strings.Builder

	mime := mimetype.Detect(data)
	mimeType := mime.String()
	content := string(data)

	if filepath == "" {
		filepath = "standard input"
	}

	// Append MIME type information if required.
	if cr.Config.MimeType {
		mimeType := fmt.Sprintf("%s (%s)\n", filepath, mimeType)

		if cr.Config.HTML || cr.Config.Markdown {
			strBuilder.WriteString("<!--\n" + mimeType + "-->\n")
		} else {
			strBuilder.WriteString(mimeType)
		}
	}

	// Format the content based on the configuration.
	switch {
	case cr.Config.HTML:
		strBuilder.WriteString(fmt.Sprintf("<code>\n%s\n</code>", content))
	case cr.Config.Markdown:
		strBuilder.WriteString(fmt.Sprintf("```\n%s\n```", content))
	case cr.Config.LineNumbers:
		lines := strings.Split(content, "\n")

		for index, line := range lines {
			strBuilder.WriteString(fmt.Sprintf("%d: %s\n", index+1, line))
		}
	default:
		strBuilder.WriteString(content)
	}

	return strBuilder.String(), nil
}
