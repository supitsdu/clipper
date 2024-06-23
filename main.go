package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

const version = "1.4.0"

// copyToClipboard writes the content to the clipboard
func copyToClipboard(contentStr string) error {
	return clipboard.WriteAll(contentStr)
}

// readFromStdin reads content from stdin
func readFromStdin() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("error reading from stdin: %v", err)
	}
	return string(input), nil
}

// readFromFile reads content from the specified file
func readFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %v", filePath, err)
	}
	return string(content), nil
}

func main() {
	// Define flags
	directText := flag.String("c", "", "Copy text directly from command line argument")
	showVersion := flag.Bool("v", false, "Show the current version of the clipper tool")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Clipper is a lightweight command-line tool for copying contents to the clipboard.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  clipper [arguments] [file ...]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nArguments:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -c <string>    Copy text directly from command line argument\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -v             Show the current version of the clipper tool\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nIf no file or text is provided, reads from standard input.\n")
	}

	flag.Parse()

	// Check if the version flag is set
	if *showVersion {
		fmt.Printf("Clipper %s\n", version)
		return
	}

	var contentStr string
	var err error

	if *directText != "" {
		// Use the provided direct text
		contentStr = *directText
	} else {
		if len(flag.Args()) > 0 {
			// Read the content from all provided file paths
			var sb strings.Builder
			for _, filePath := range flag.Args() {
				content, err := readFromFile(filePath)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				sb.WriteString(content + "\n")
			}
			contentStr = sb.String()
		} else {
			// Read from stdin
			contentStr, err = readFromStdin()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Write the content to the clipboard
	err = copyToClipboard(contentStr)
	if err != nil {
		fmt.Printf("Error copying content to clipboard: %v\n", err)
		os.Exit(1)
	}

	// Print success message
	fmt.Println("Clipboard updated successfully. Ready to paste!")
}
