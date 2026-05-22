package main

import (
	"fmt"
	"os"
)

func finish(output string, code int) {
	if output != "" {
		fmt.Print(output)
	}
	os.Exit(code)
}

func main() {
	p, err := NewPicker()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	r := p.Run()
	finish(r.Output, r.Code)
}
