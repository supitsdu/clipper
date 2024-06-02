package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	// Define flags
	directText := flag.String("c", "", "Copy text directly from command line argument")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Clipper is a lightweight command-line tool for copying contents to the clipboard.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  clipper [arguments]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nArguments:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -c <string>    Copy text directly from command line argument\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nIf no file or text is provided, reads from standard input.\n")
	}

	flag.Parse()

	var contentStr string

	if *directText != "" {
		// Use the provided direct text
		contentStr = *directText
	} else if len(flag.Args()) == 1 {
		// Read the content from the file path provided
		filePath := flag.Arg(0)
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file '%s': %v\n", filePath, err)
			os.Exit(1)
		}
		contentStr = string(content)
	} else {
		// Read from stdin
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		contentStr = string(input)
	}

	// Write the content to the clipboard
	err := clipboard.WriteAll(contentStr)
	if err != nil {
		fmt.Printf("Error copying content to clipboard: %v\n", err)
		os.Exit(1)
	}

	// Print success message
	fmt.Println("Clipboard updated successfully. Ready to paste!")
}
