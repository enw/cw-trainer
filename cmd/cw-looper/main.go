//go:build darwin

package main

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"cw-looper/internal/terminal"
	"flag"
	"fmt"
	"os"
)

func run() error {
	path, err := library.DefaultPath()
	if err != nil {
		return err
	}
	dataPath := flag.String("data", path, "library JSON path (also locks this library)")
	flag.Parse()
	store, err := library.Open(*dataPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", *dataPath, err)
	}
	defer store.Close()
	data, err := store.Load()
	if err != nil {
		return fmt.Errorf("load %s: %w; existing data preserved, correct the file before retrying", *dataPath, err)
	}
	runner := output.NewRunner(&output.Mac{})
	defer runner.Close()
	_, err = tea.NewProgram(terminal.New(data, runner, store.Save)).Run()
	return err
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "cw-looper:", err)
		os.Exit(1)
	}
}
