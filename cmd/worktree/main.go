package main

import (
	"fmt"
	"os"

	"github.com/ridge/worktree"
)

func main() {
	version, err := worktree.CurrentVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get worktree version: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(version)
}
