package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scnewma/lpass-auditor/command"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s /path/to/csv\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	command := command.Audit{}

	err := command.Execute(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
