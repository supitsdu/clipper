package options

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Text         string
	FilePaths    []string
	HTML         bool
	Markdown     bool
	MimeType     bool
	LineNumbers  bool
	ShouldFormat bool
	ShowVersion  bool
}

// Package-level variables for version information (set at build time or default).
var (
	Version       = "dev"        // Default for development builds.
	BuildMetadata = "git/source" // Default for development builds.
)

// GetVersion formats the version string for display.
func GetVersion() string {
	versionStr := strings.TrimSpace(Version)

	if BuildMetadata != "" {
		versionStr += " " + strings.TrimSpace(BuildMetadata)
	}

	return versionStr
}

// ParseFlags parses the command-line flags and arguments.
func ParseFlags() *Config {
	text := flag.String("c", "", "Copy text directly from command line argument")
	htmlWrap := flag.Bool("Html", false, "Each file data is put within an HTML5 codeblock.")
	markdownWrap := flag.Bool("Markdown", false, "Each file data is put within an Markdown codeblock.")
	mimeType := flag.Bool("Mime", false, "Include mimetype for each file")
	lineNumbers := flag.Bool("LineNumbers", false, "Add line numbers to the content.")
	showVersion := flag.Bool("v", false, "Show the current version of the clipper tool")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Clipper is a lightweight command-line tool for copying contents to the clipboard.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  clipper [arguments] [file ...]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nArguments:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nIf no file or text is provided, reads from standard input.\n")
	}

	flag.CommandLine.SetOutput(os.Stderr)

	flag.Parse()

	return &Config{
		Text:         *text,
		FilePaths:    flag.Args(),
		HTML:         *htmlWrap,
		Markdown:     *markdownWrap,
		MimeType:     *mimeType,
		LineNumbers:  *lineNumbers,
		ShowVersion:  *showVersion,
		ShouldFormat: *htmlWrap || *mimeType || *markdownWrap || *lineNumbers,
	}
}
