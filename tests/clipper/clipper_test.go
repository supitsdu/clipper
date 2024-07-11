package clipper_test

import (
	"testing"

	"github.com/atotto/clipboard"
	"github.com/supitsdu/clipper/cli/clipper"
	"github.com/supitsdu/clipper/tests"
)

func TestClipboardWriter(t *testing.T) {
	t.Run("DefaultClipboardWriter", func(t *testing.T) {
		if tests.IsCIEnvironment() {
			t.Skip("Skipping clipboard test in Github Actions. Helps avoid errors when on CI environments.")
		}

		// Create a DefaultClipboardWriter
		writer := clipper.DefaultClipboardWriter{}

		// Write some content to the clipboard
		err := writer.Write(tests.SampleText32)
		if err != nil {
			t.Errorf("Error writing to clipboard: %v", err)
		}

		// Check the clipboard content
		clipboardContent, err := clipboard.ReadAll()
		if err != nil {
			t.Errorf("Error reading from clipboard: %v", err)
		}

		// Check the content
		if clipboardContent != tests.SampleText32 {
			t.Errorf("Expected '%s', got '%s'", tests.SampleText32, clipboardContent)
		}
	})
}
