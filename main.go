package main

import (
	"fmt"
	"os"
)

func main() {
	cli := NewCLI(os.Stdout, os.Stderr, os.Stdin)
	err := cli.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
