package main

import (
	"fmt"
	"os"

	"github.com/supitsdu/clipper/cli/clipper"
	"github.com/supitsdu/clipper/cli/options"
)

func main() {
	config := options.ParseFlags()
	writer := clipper.DefaultClipboardWriter{}

	msg, err := clipper.Run(config, writer)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}

	if msg != "" {
		fmt.Printf("Clipper %s\n", msg)
		os.Exit(0)
	}
}
