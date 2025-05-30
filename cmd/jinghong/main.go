package main

import (
	"jonghong/internal/jinghong"
	"os"
)

func main() {
	cmd := jinghong.NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
