package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	// Define the flag for copying text directly from command line argument
	directText := flag.String("c", "", "Copy text directly from command line argument")
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
