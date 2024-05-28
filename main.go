package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "github.com/atotto/clipboard" // Import the clipboard library
)

func main() {
    // Check if the correct number of arguments are provided
    if len(os.Args) != 2 {
        fmt.Println("Usage: clipper <file_path>")
        os.Exit(1)
    }

    // Get the file path from command line arguments
    filePath := os.Args[1]

    // Read the content of the file
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    // Write the content to the clipboard
    err = clipboard.WriteAll(string(content))
    if err != nil {
        fmt.Printf("Error copying content to clipboard: %v\n", err)
        os.Exit(1)
    }

    // Print success message
    fmt.Println("Clipboard updated with file content. Ready to paste!")
}
