package main

import (
	"fmt"
	"os"

	"github.com/tenntenn/lualin/lualin"
	"github.com/yuin/gopher-lua/parse"
)

func main() {
	chunk, err := parse.Parse(os.Stdin, "sample")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}

	if err := lualin.Lint(os.Stdout, chunk); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
