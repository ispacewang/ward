package main

import (
	"os"

	"github.com/example/docgen/internal/app"
)

func main() {
	if err := app.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
