package options

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	DirectText  string
	ShowVersion bool
	Args        []string
}

// Package-level variables for version information (set at build time or default)
var (
	Version       = "dev"        // Default for development builds
	BuildMetadata = "git/source" // Default for development builds
)

// GetVersion formats the version string for display
func GetVersion() string {
	if BuildMetadata != "" {
		return strings.TrimSpace(Version) + " " + strings.TrimSpace(BuildMetadata)
	} else {
		return strings.TrimSpace(Version)
	}
}

// ParseFlags parses the command-line flags and arguments.
func ParseFlags() *Config {
	directText := flag.String("c", "", "Copy text directly from command line argument")
	showVersion := flag.Bool("v", false, "Show the current version of the clipper tool")

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

	return &Config{
		DirectText:  *directText,
		ShowVersion: *showVersion,
		Args:        flag.Args(),
	}
}
