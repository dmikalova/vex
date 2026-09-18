package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// reviewStateFile records which card files the human has already reviewed this
// cycle, so `review` never surfaces the same file again until the whole pool has
// been seen. It is personal progress rather than shared state, so it is
// gitignored; it lives at the repo root because mage runs from there.
const reviewStateFile = ".card-review.json"

// defaultReviewBatch is how many card files a review run opens when -n is unset.
const defaultReviewBatch = 10

// reviewState is the persisted set of already-reviewed card files, stored as
// repo-root-relative slash paths.
type reviewState struct {
	Viewed []string `json:"viewed"`
}

// review opens a random batch of not-yet-reviewed card files in VS Code and
// records them as reviewed, cycling: once every eligible file has been seen the
// slate clears and a fresh pass begins. Eligible files are the card definitions
// under internal/cards/sets, minus test files, the generated 0set.go catalogs,
// and build-excluded (`//go:build todo`) stubs — the files with a real ability to
// read. Run `mage tool:review` (or `mage tool:review -n=5`).
func review(args []string) error {
	batch := defaultReviewBatch
	for _, a := range args {
		rest, ok := strings.CutPrefix(a, "-n=")
		if !ok {
			return fmt.Errorf("usage: cardlookup review [-n=<count>]")
		}
		n, err := strconv.Atoi(rest)
		if err != nil || n <= 0 {
			return fmt.Errorf("cardlookup review: -n must be a positive integer, got %q", rest)
		}
		batch = n
	}

	files, err := reviewableCardFiles()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("cardlookup review: no card files found under %s", setsDir)
	}

	state, err := loadReviewState()
	if err != nil {
		return err
	}
	// Forget remembered files that no longer exist, so "all reviewed" stays true.
	viewed := pruneMissing(state.Viewed, files)

	remaining := subtractPaths(files, viewed)
	if len(remaining) == 0 {
		fmt.Printf("All %d card files reviewed — starting a new cycle.\n", len(files))
		viewed = nil
		remaining = append(remaining, files...)
	}

	rand.Shuffle(len(remaining), func(i, j int) {
		remaining[i], remaining[j] = remaining[j], remaining[i]
	})
	if batch > len(remaining) {
		batch = len(remaining)
	}
	pick := remaining[:batch]
	sort.Strings(pick)

	viewed = append(viewed, pick...)
	sort.Strings(viewed)
	if err := saveReviewState(reviewState{Viewed: viewed}); err != nil {
		return err
	}

	fmt.Printf(
		"Opening %d card file(s) for review (%d of %d seen this cycle):\n",
		len(pick), len(viewed), len(files),
	)
	for _, p := range pick {
		fmt.Printf("  %s\n", p)
	}
	return openInEditor(pick)
}

// reviewableCardFiles lists the card definition files under internal/cards/sets
// worth reviewing: every .go file except test files, the generated 0set.go
// catalogs, and build-excluded stubs. Paths are repo-root-relative slash paths.
func reviewableCardFiles() ([]string, error) {
	var files []string
	err := filepath.WalkDir(setsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, "_test.go") ||
			name == "0set.go" {
			return nil
		}
		if hasBuildTodo(path) {
			return nil
		}
		files = append(files, filepath.ToSlash(path))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cardlookup review: %w", err)
	}
	sort.Strings(files)
	return files, nil
}

// loadReviewState reads the reviewed-files record, treating a missing file as an
// empty slate.
func loadReviewState() (reviewState, error) {
	data, err := os.ReadFile(reviewStateFile)
	if os.IsNotExist(err) {
		return reviewState{}, nil
	}
	if err != nil {
		return reviewState{}, fmt.Errorf("cardlookup review: reading %s: %w", reviewStateFile, err)
	}
	var s reviewState
	if err := json.Unmarshal(data, &s); err != nil {
		return reviewState{}, fmt.Errorf("cardlookup review: parsing %s: %w", reviewStateFile, err)
	}
	return s, nil
}

// saveReviewState writes the reviewed-files record.
func saveReviewState(s reviewState) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("cardlookup review: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(reviewStateFile, data, 0o644); err != nil {
		return fmt.Errorf("cardlookup review: writing %s: %w", reviewStateFile, err)
	}
	return nil
}

// pruneMissing drops remembered paths that are no longer among the current files.
func pruneMissing(viewed, files []string) []string {
	present := make(map[string]bool, len(files))
	for _, f := range files {
		present[f] = true
	}
	var kept []string
	for _, v := range viewed {
		if present[v] {
			kept = append(kept, v)
		}
	}
	return kept
}

// subtractPaths returns the files not present in viewed.
func subtractPaths(files, viewed []string) []string {
	seen := make(map[string]bool, len(viewed))
	for _, v := range viewed {
		seen[v] = true
	}
	var out []string
	for _, f := range files {
		if !seen[f] {
			out = append(out, f)
		}
	}
	return out
}

// openInEditor opens the files in VS Code, reusing the current window. A missing
// `code` CLI is not an error — the paths were already printed for the human to
// open by hand.
func openInEditor(files []string) error {
	bin, err := exec.LookPath("code")
	if err != nil {
		fmt.Println("(the `code` CLI is not on PATH — open the files above by hand)")
		return nil
	}
	cmd := exec.Command(bin, append([]string{"--reuse-window"}, files...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
