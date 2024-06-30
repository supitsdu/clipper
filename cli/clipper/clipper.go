package clipper

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/supitsdu/clipper/cli/options"
	"github.com/supitsdu/clipper/cli/reader"
)

// ClipboardWriter defines an interface for writing content to the clipboard.
type ClipboardWriter interface {
	Write(content string) error
}

// DefaultClipboardWriter writes content to the clipboard using the default clipboard implementation.
type DefaultClipboardWriter struct{}

// Write writes the given content to the clipboard.
func (c DefaultClipboardWriter) Write(content string) error {
	return clipboard.WriteAll(content)
}

// Run executes the clipper tool logic based on the provided configuration.
func Run(config *options.Config, writer ClipboardWriter) (string, error) {
	if *config.ShowVersion {
		return options.Version, nil
	}

	readers := reader.GetReaders(config.Args)

	// Aggregate the content from the provided sources.
	content, err := reader.ParseContent(config.DirectText, readers...)
	if err != nil {
		return "", fmt.Errorf("parsing content: %w", err)
	}

	// Write the parsed content to the provided clipboard.
	if err = writer.Write(content); err != nil {
		return "", fmt.Errorf("copying content to clipboard: %w", err)
	}

	return "updated clipboard successfully. Ready to paste!", nil
}
