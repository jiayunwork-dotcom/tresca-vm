package main

import (
	"os"

	"tresca-vm/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
