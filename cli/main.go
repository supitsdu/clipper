package main

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"

	"github.com/supitsdu/clipper/cli/clipper"
	"github.com/supitsdu/clipper/cli/options"
	"github.com/supitsdu/clipper/cli/reader"
)

func main() {
	config := options.ParseFlags() // Parse command-line flags

	if config.ShowVersion {
		fmt.Printf("Clipper %s\n", options.GetVersion())
		os.Exit(0)
	}

	writer := clipper.DefaultClipboardWriter{} // Create the default clipboard writer

	reader := reader.ContentReader{Config: config}

	bytesLength, err := clipper.Run(reader, writer) // Run the main Clipper logic
	if err != nil {
		fmt.Printf("Error %s.\n", err)
		os.Exit(1)
	}

	fmt.Printf("Copied %s to the clipboard. Ready to paste!\n", humanize.Bytes(uint64(bytesLength)))
}
