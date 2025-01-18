package main

import (
	"fmt"
	"os"
)

var version = "v0.0.1"

func run() error {
	if err := startHTTPServer(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func main() {
	// versionという引数が指定された場合はバージョン情報を表示する
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Fprintf(os.Stdout, "version: %s\n", version)
		os.Exit(0)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
