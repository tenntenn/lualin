package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/shiena/ansicolor"
	"github.com/tenntenn/lualin/lualin"
	"github.com/wsxiaoys/terminal/color"
	"github.com/yuin/gopher-lua/parse"
)

var (
	stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
)

func printError(fn string, err *lualin.LintError) {
	switch err.Rule.Level() {
	case lualin.Error:
		color.Fprintf(stderr, "@{r}%s:%s\n", fn, err.Error())
	case lualin.Warning:
		color.Fprintf(stdout, "@{y}%s:%s\n", fn, err.Error())
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "lualin"
	app.Usage = "lua lint"
	app.Action = func(c *cli.Context) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					color.Fprintf(stderr, "@{r}Error: %s\n", err.Error())
				}
			}
		}()

		for _, fn := range c.Args() {
			f, err := os.Open(fn)
			if err != nil {
				panic(err)
			}

			chunk, err := parse.Parse(f, fn)
			if err != nil {
				panic(err)
			}

			if err := lualin.Lint(os.Stdout, chunk); err != nil {
				switch err.(type) {
				case *lualin.LintError:
					linterr, _ := err.(*lualin.LintError)
					printError(fn, linterr)
				case lualin.LintErrors:
					errs, _ := err.(lualin.LintErrors)
					for _, linterr := range errs {
						printError(fn, linterr)
					}
				default:
					panic(err)
				}
			}
		}
	}

	app.Run(os.Args)
}
