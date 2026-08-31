//go:build mage

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
)

// Cover runs tests and reports engine coverage, failing if it drops below 100%.
func Cover() error {
	if err := sh.RunV(
		"go",
		"test",
		"./internal/engine/",
		"-coverprofile=coverage.out",
	); err != nil {
		return err
	}
	out, err := sh.Output("go", "tool", "cover", "-func=coverage.out")
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	total := lines[len(lines)-1]
	fmt.Println(total)
	fields := strings.Fields(total)
	pct, err := strconv.ParseFloat(strings.TrimSuffix(fields[len(fields)-1], "%"), 64)
	if err != nil {
		return fmt.Errorf("parsing coverage percent from %q: %w", total, err)
	}
	if pct < 100 {
		fmt.Println("Functions below 100%:")
		for _, l := range lines[:len(lines)-1] {
			if !strings.HasSuffix(l, "100.0%") {
				fmt.Println("  " + l)
			}
		}
		return fmt.Errorf("engine coverage %.1f%% is below 100%%", pct)
	}
	return nil
}
