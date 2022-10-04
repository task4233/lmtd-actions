package main

import (
	"context"
	"fmt"
	"os"

	"github.com/task4233/lmtd-cheker"
)

// build時のldflagsでembed
// -ldflags "-X github.com/task4233/lmtd-cheker/cmd/lmtd-cheker/main.version={version}"
var version string

func main() {
	lmtd := lmtd.LMTd{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := lmtd.Run(context.Background(), version, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "failed lmtd.Run: %s", err.Error())
		os.Exit(1)
	}
}
