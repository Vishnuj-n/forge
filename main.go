package main

import (
	"os"

	"forge/cmd/forge"
)

func main() {
	if err := forge.Execute(); err != nil {
		os.Exit(1)
	}
}
