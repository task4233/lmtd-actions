package main

import (
	"context"
	"fmt"
	"os"

	"github.com/task4233/lmtd-actions"
)

// build時のldflagsでembed
// -ldflags "-X github.com/task4233/lmtd-actions/cmd/lmtd-actions/main.version={version}"
var version string

func main() {
	lmtd := lmtd.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := lmtd.Run(context.Background(), version, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "failed lmtd.Run: %s", err.Error())
		os.Exit(1)
	}
}
