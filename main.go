package main

import (
	"fmt"
	"os"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainE() error {
	fmt.Println("Hello, World") //nolint

	return nil
}
