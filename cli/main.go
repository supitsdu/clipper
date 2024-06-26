package main

import (
	"github.com/supitsdu/clipper/cli/clipper"
	"github.com/supitsdu/clipper/cli/options"
)

func main() {
	config := options.ParseFlags()
	clipper.Run(config)
}
