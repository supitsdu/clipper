package main

import (
	"fmt"
	"os"

	"github.com/supitsdu/clipper/cli/clipper"
	"github.com/supitsdu/clipper/cli/options"
)

func main() {
	config := options.ParseFlags() // Parse command-line flags

	if *config.ShowVersion {
		fmt.Printf("Clipper %s\n", options.GetVersion())
		os.Exit(0)
	}

	writer := clipper.DefaultClipboardWriter{} // Create the default clipboard writer

	msg, err := clipper.Run(config, writer) // Run the main Clipper logic
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1) // Exit with an error code
	}

	fmt.Println(msg)
}
