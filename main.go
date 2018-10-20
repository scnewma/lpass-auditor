package main

import (
	"fmt"
	"os"

	"github.com/scnewma/auditor/command"
)

func main() {
	command := command.Audit{}

	err := command.Execute(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
