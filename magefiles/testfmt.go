//go:build mage

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/tabwriter"
)

// modulePath is stripped from every package path in the test report so a line is
// the short package (internal/engine) rather than the full import path, which is
// what makes the report fit an 80-column terminal without wrapping.
const modulePath = "github.com/dmikalova/vactrol/"

// goTest runs `go test` with the given args and prints a tidied, column-aligned
// report in place of go test's ragged, path-heavy output: the module prefix is
// stripped, the coverage suffix is shortened, and the status/package/detail
// columns line up whatever the package names' lengths. Lines go test emits that
// are not per-package summaries (failure details, panics, build errors) pass
// through verbatim, so a tidy report never hides why a run went red. It returns
// the run's error unchanged.
func goTest(args ...string) error {
	cmd := exec.Command("go", append([]string{"test"}, args...)...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	runErr := cmd.Run()
	fmt.Print(tidyTestOutput(buf.String()))
	return runErr
}

// tidyTestOutput rewrites go test's summary lines into an aligned table and
// leaves every other line as it is. A summary line is one that starts with `ok`,
// `?`, or `FAIL` followed by a package; anything else (a failing test's own
// output) is flushed through in place so its order and detail survive.
func tidyTestOutput(raw string) string {
	var out bytes.Buffer
	tw := tabwriter.NewWriter(&out, 0, 0, 2, ' ', 0)
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Split(line, "\t")
		switch {
		case strings.HasPrefix(line, "ok") && len(fields) >= 2:
			detail := ""
			if len(fields) >= 3 {
				detail = strings.TrimSpace(fields[2])
			}
			cov := ""
			if len(fields) >= 4 {
				cov = shortCoverage(fields[3])
			}
			fmt.Fprintf(tw, "ok\t%s\t%s\t%s\n", shortPkg(fields[1]), detail, cov)
		case strings.HasPrefix(line, "?") && len(fields) >= 2:
			fmt.Fprintf(tw, "?\t%s\t%s\t\n", shortPkg(fields[1]), "no test files")
		case strings.HasPrefix(line, "FAIL\t") && len(fields) >= 2:
			detail := ""
			if len(fields) >= 3 {
				detail = strings.TrimSpace(fields[2])
			}
			fmt.Fprintf(tw, "FAIL\t%s\t%s\t\n", shortPkg(fields[1]), detail)
		default:
			// A non-summary line breaks the aligned block, so the table so far is
			// flushed before it prints, keeping output in the order go test emitted it.
			tw.Flush()
			fmt.Fprintln(&out, line)
		}
	}
	tw.Flush()
	// tabwriter pads the detail column to its width even when the coverage column
	// after it is empty, so each padded line is right-trimmed to drop that tail.
	var trimmed strings.Builder
	sc = bufio.NewScanner(strings.NewReader(out.String()))
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		trimmed.WriteString(strings.TrimRight(sc.Text(), " "))
		trimmed.WriteByte('\n')
	}
	return trimmed.String()
}

// shortPkg drops the module prefix from a package path so the report names it the
// way the tree does, e.g. internal/engine.
func shortPkg(field string) string {
	return strings.TrimPrefix(strings.TrimSpace(field), modulePath)
}

// shortCoverage trims go test's "coverage: X% of statements in <pkgs>" down to
// "cov X%": the trailing "of statements in ..." is the noise that makes the line
// wrap, and which package the count is for is already the row's package.
func shortCoverage(field string) string {
	s := strings.TrimPrefix(strings.TrimSpace(field), "coverage: ")
	if i := strings.Index(s, " of statements"); i >= 0 {
		s = s[:i]
	}
	if s == "" {
		return ""
	}
	return "cov " + s
}
