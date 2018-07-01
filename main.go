package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"go.felesatra.moe/subcommands"
)

var progName = os.Args[0]
var commands = make([]subcommands.Cmd, 0, 5)

func main() {
	if err := subcommands.Run(commands, os.Args[1:]); err != nil {
		fmt.Fprint(os.Stderr, err)
		usage(os.Stderr)
		os.Exit(1)
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Valid commands:")
	for _, c := range commands {
		fmt.Fprintln(w, c.Name())
	}
}

func errCmd(f func([]string) error) func([]string) {
	return func(args []string) {
		if err := f(args); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
