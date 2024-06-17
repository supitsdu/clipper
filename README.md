# Clipper - Ready to Paste

Clipper is a lightweight command-line tool written in Go for copying contents to the clipboard. Quickly and easily copy the contents of any file, command output, or even directly pass text to your clipboard inside your favorite terminal, streamlining your workflow and saving you time.

## Features ✨

- **Cross-Platform Compatibility:** Clipper works seamlessly on **Linux**, **macOS**, and **Windows**, providing consistent clipboard functionality across different operating systems.
- **Simple Usage:** With a straightforward command-line interface, Clipper makes it easy to copy file contents to the clipboard with just a single command.
- **Fast and Efficient:** Clipper is designed for performance and efficiency, allowing you to copy contents to the clipboard quickly and without unnecessary overhead.
- **Convenient:** Directly copy text, command output, or file contents to the clipboard in a consistent way across all supported platforms.

## Installation

### Using the Install Script

You can install Clipper using the following command:

```sh
curl -sSL https://raw.githubusercontent.com/supitsdu/clipper/main/install.sh | sh -c
```

Or:

```sh
wget -qO- https://raw.githubusercontent.com/supitsdu/clipper/main/install.sh | sh
```

> Note: If you prefer not to use `curl ... | sh`, you can manually download and install Clipper as described below.

### Manual Installation

To use Clipper, download the appropriate binary for your operating system from the [releases page](https://github.com/supitsdu/clipper/releases) and place it in your desired location. Add the location of the binary to your system's PATH environment variable to access Clipper from anywhere on your system.

## Usage

Copy the contents of a file to the clipboard:

```sh
clipper ~/.ssh/id_ed25519.pub
```

Copy the output of any command:

```sh
go version | clipper
```

Directly copy text:

```sh
clipper -c "My awesome content!"
```

## Contributing 🤝

Contributions to Clipper are welcome! Here are a few ways you can contribute:

- **Report Bugs:** If you encounter any bugs or unexpected behavior while using Clipper, please [open an issue](https://github.com/supitsdu/clipper/issues) on GitHub to report it.
- **Request Features:** Have an idea for a new feature or improvement? [Open an issue](https://github.com/supitsdu/clipper/issues) to share your suggestion and start a discussion.
- **Submit Pull Requests:** If you're comfortable with Go programming, you can contribute directly to the development of Clipper by submitting pull requests. Fork the repository, make your changes, and submit a pull request for review.

## Building from Source

To build Clipper from source, you'll need to have Go installed on your system, as well as the Make tool.

### Requirements

- **Go Lang:** Clipper is written in Go, so you'll need to have Go installed on your system. You can download and install it from the [official Go website](https://golang.org/).
- **Make Tool:** Building Clipper from source requires the Make tool to automate the build process. Make is commonly pre-installed on Unix-like systems, but you may need to install it manually on Windows.

Once you have the necessary requirements installed, clone the repository and run the following command in the project directory:

```sh
make
```

This command will build binaries for Windows, Linux, and macOS in the `bin` directory.

## License

Clipper is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute it for any purpose.
