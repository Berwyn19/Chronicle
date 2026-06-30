package main

import (
	"fmt"
	"os"

	"github.com/berwyn/chronicle/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "chronicle:", err)
		os.Exit(1)
	}
}
