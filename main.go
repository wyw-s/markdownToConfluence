package main

import (
	"markdownToConfluence/cmd"
)

// version of markdown2confluence. Overwritten during build
var version = "0.0.0"

func main() {
	cmd.Execute(version)
}
