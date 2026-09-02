package sim

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// CorpusDir is FuzzPlay's seed corpus, relative to the sim package. Go replays
// every file in it as a subtest of a plain `go test`, so a script saved here is a
// permanent regression, not just fuzzer scratch.
const CorpusDir = "testdata/fuzz/FuzzPlay"

// seedLen is the prefix of a script the decoder reads as the game's seed. It
// picks the decks, so minimizing must never trim into it.
const seedLen = 8

// incidental folds the parts of a violation message that vary between scripts
// tripping the same bug: the numbers (turn, step, card id, power totals) and any
// dumped struct, whose booleans differ run to run. The card name is left alone —
// the same invariant broken on a different card is worth its own corpus entry.
var incidental = regexp.MustCompile(`[0-9]+|\{[^{}]*\}`)

// failureKey names the bug behind a violation: the message with its incidental
// detail folded away, hashed. Two scripts that fail the same way share a key, and
// so share one corpus file.
func failureKey(err error) string {
	sum := sha256.Sum256([]byte(incidental.ReplaceAllString(err.Error(), "#")))
	return hex.EncodeToString(sum[:])
}

// Minimize trims a failing script down to the shortest prefix that still fails
// the same way. A soak script is mostly padding the decoder never reads, so this
// usually turns kilobytes of random bytes into a script short enough to read.
func Minimize(script []byte) []byte {
	key, failing := failureOf(script)
	if !failing {
		return script
	}
	best := script
	for step := len(best) / 2; step > 0; step /= 2 {
		for len(best)-step > seedLen {
			candidate := best[:len(best)-step]
			if k, ok := failureOf(candidate); !ok || k != key {
				break
			}
			best = candidate
		}
	}
	return best
}

// failureOf reports the key of the violation a script trips, if it trips one.
func failureOf(script []byte) (string, bool) {
	err := Simulate(script)
	if err == nil {
		return "", false
	}
	return failureKey(err), true
}

// SaveCorpus files a failing script as a regression, minimized and named for the
// bug it catches rather than for its own bytes, so one bug keeps one entry
// however many scripts find it.
func SaveCorpus(dir string, script []byte, err error) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := filepath.Join(dir, failureKey(err))
	body := fmt.Sprintf("go test fuzz v1\n[]byte(%q)\n", Minimize(script))
	if writeErr := os.WriteFile(name, []byte(body), 0o644); writeErr != nil {
		return "", writeErr
	}
	return name, nil
}

// PruneReport is the outcome of a corpus prune, in the terms a caller reports.
type PruneReport struct {
	Scanned  int
	Kept     int      // entries still reproducing a bug, one per distinct failure
	Fixed    int      // entries whose bug is fixed, so they no longer fail
	Failures []string // the distinct violations still in the corpus
}

// PruneCorpus replays every corpus entry and rewrites the directory as one
// minimized entry per bug that still reproduces. Entries whose bug is fixed no
// longer fail, so they no longer pin anything down and are dropped — the corpus
// is the list of open findings, not an archive of every script a soak ever saw.
func PruneCorpus(dir string) (PruneReport, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return PruneReport{}, err
	}

	var report PruneReport
	keep := map[string][]byte{} // failure key -> the shortest script that trips it
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		script, readErr := readCorpusFile(path)
		if readErr != nil {
			return report, fmt.Errorf("%s: %w", entry.Name(), readErr)
		}
		report.Scanned++

		failure := Simulate(script)
		if failure == nil {
			report.Fixed++
			if rmErr := os.Remove(path); rmErr != nil {
				return report, rmErr
			}
			continue
		}
		key := failureKey(failure)
		best, seen := keep[key]
		if !seen {
			report.Failures = append(report.Failures, failure.Error())
		}
		if minimized := Minimize(script); !seen || len(minimized) < len(best) {
			keep[key] = minimized
		}
		if entry.Name() != key {
			if rmErr := os.Remove(path); rmErr != nil {
				return report, rmErr
			}
		}
	}

	for key, script := range keep {
		body := fmt.Sprintf("go test fuzz v1\n[]byte(%q)\n", script)
		if writeErr := os.WriteFile(filepath.Join(dir, key), []byte(body), 0o644); writeErr != nil {
			return report, writeErr
		}
	}
	report.Kept = len(keep)
	return report, nil
}

// readCorpusFile decodes the script out of a Go fuzz corpus file, whose body is a
// version line followed by one Go-quoted []byte literal.
func readCorpusFile(path string) ([]byte, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	for line := range strings.Lines(string(body)) {
		quoted, ok := strings.CutPrefix(strings.TrimSpace(line), "[]byte(")
		if !ok {
			continue
		}
		quoted, ok = strings.CutSuffix(quoted, ")")
		if !ok {
			continue
		}
		unquoted, unquoteErr := strconv.Unquote(quoted)
		if unquoteErr != nil {
			return nil, unquoteErr
		}
		return []byte(unquoted), nil
	}
	return nil, errors.New("no []byte literal in corpus file")
}
