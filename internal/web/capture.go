package web

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"sync"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/cards"
	"github.com/dmikalova/vex/internal/engine"
)

// This file is the other half of the quarantine in game_persist.go. Quarantine
// keeps a failing match where only that one browser can see it; a capture carries
// it out to a file an agent who was never at the keyboard can replay.
//
// The client is wasm in a browser sandbox and cannot write a file, so a capture is
// a round trip: the client POSTs the failing record to the local dev server, which
// writes it under CaptureDir. The endpoint exists only when DevEnabled, so a
// deployed build has no handler and the POST simply fails — and a failed POST is
// silent, because the quarantine slot still holds the reproduction locally. The
// two halves degrade into each other.
//
// A capture on disk is an OPEN FINDING, not an archive entry: TestCaptures
// replays every one and fails, so a capture sitting in a commit is a lapse.
// Findings are fixed and pruned (`mage capturePrune`) before commit, exactly as
// the FuzzPlay seed corpus is.

// CapturePath is the dev-only endpoint a replay failure posts its record to.
const CapturePath = "/debug/capture"

// CaptureDir is where captures are written, relative to the repo root — the
// working directory `mage web` runs the dev server from. They are committed, not
// ignored, which is what makes ci:check surface one to somebody else.
const CaptureDir = "internal/web/testdata/capture"

// captureTestDir is the same directory as this package's own tests see it, since
// `go test` runs with the package directory as the working directory.
const captureTestDir = "testdata/capture"

// Capture is one recorded replay failure: the match that would not replay, the
// panic it died on, and the stamps that say what tree it was recorded against.
//
// The stamps are the whole point of the format. A capture is only evidence
// against the tree that produced it: Version moves when the command log's meaning
// changes, and PoolDigest moves when the implemented card pool does — which
// happens constantly, without a version bump, and makes a capture recorded against
// another pool impossible to debug. A capture whose stamps do not match the
// current tree is skipped with a note, never failed (see Replay).
type Capture struct {
	// Version is snapshotVersion at the time of recording.
	Version int
	// Commit is the vcs revision the dev server was built from, recorded for the
	// agent investigating the finding. It is deliberately NOT used to decide
	// staleness: an ancestry test would be wrong across branches.
	Commit string
	// PoolDigest is a hash over the implemented card database (see poolDigest).
	PoolDigest string
	// Panic is the recovered panic value, rendered.
	Panic string
	// Stack is the goroutine stack at the recover.
	Stack string
	// Snapshot is the match itself: seed, sets, and the command log that fails.
	Snapshot snapshot

	// Path is where this capture was loaded from, so a prune can delete it. It is
	// not part of the file.
	Path string `json:"-"`
}

// CaptureStatus is what replaying one capture found.
type CaptureStatus int

const (
	// CaptureStale means the capture was recorded against a different command-log
	// version or card pool, so replaying it here would prove nothing.
	CaptureStale CaptureStatus = iota
	// CaptureFixed means the capture replays cleanly now: the finding is closed and
	// the file can be pruned.
	CaptureFixed
	// CaptureReproduced means the capture still fails, the same way or another.
	CaptureReproduced
)

// Replay re-deals the captured match and feeds its command log back, reporting
// what happened and a one-line detail for the reader. A stale capture is reported
// stale without being replayed at all — it is not evidence, and must never fail a
// gate. Proven by TestAStaleCaptureIsSkipped and TestAFailingCaptureReproduces.
func (c Capture) Replay() (status CaptureStatus, detail string) {
	if c.Version != snapshotVersion {
		return CaptureStale, fmt.Sprintf(
			"recorded under command-log version %d; this tree writes %d",
			c.Version, snapshotVersion)
	}
	if c.PoolDigest != poolDigest() {
		return CaptureStale, "recorded against a different card pool"
	}
	// The failure being replayed is, by construction, one that panics or diverges;
	// catching it here is what turns it into a result rather than a dead test binary.
	defer func() {
		if r := recover(); r != nil {
			status, detail = CaptureReproduced, fmt.Sprint(r)
		}
	}()
	g := &game{}
	if err := g.replayRecord(c.Snapshot.Record); err != nil {
		return CaptureReproduced, err.Error()
	}
	return CaptureFixed, ""
}

// LoadCaptures reads every capture in a directory, in filename order. A missing
// directory is not an error: no findings are open.
func LoadCaptures(dir string) ([]Capture, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Capture
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var c Capture
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		c.Path = path
		out = append(out, c)
	}
	return out, nil
}

// PruneReport is what one prune pass did, for the caller to print.
type PruneReport struct {
	// Scanned is how many capture files were read.
	Scanned int
	// Fixed is how many were deleted because their fault no longer reproduces.
	Fixed int
	// Stale is how many were deleted because they were recorded against another
	// command-log version or card pool, so they can never be evidence again.
	Stale int
	// Open lists the findings that still reproduce, each with its detail. These are
	// what is left on disk.
	Open []string
}

// PruneCaptures replays every capture in dir and deletes the ones that are no
// longer open findings: the faults that have been fixed, and the entries recorded
// against a tree this one can no longer be compared with. What remains is the list
// of findings still to fix — the directory is that list, not an archive of every
// failure a dev server ever saw. Sibling of sim.PruneCorpus.
func PruneCaptures(dir string) (PruneReport, error) {
	caps, err := LoadCaptures(dir)
	if err != nil {
		return PruneReport{}, err
	}
	report := PruneReport{Scanned: len(caps)}
	for _, c := range caps {
		status, detail := c.Replay()
		if status == CaptureReproduced {
			report.Open = append(report.Open, fmt.Sprintf("%s: %s", c.Path, detail))
			continue
		}
		if err := os.Remove(c.Path); err != nil {
			return report, err
		}
		if status == CaptureStale {
			report.Stale++
		} else {
			report.Fixed++
		}
	}
	return report, nil
}

// captureName is the file a capture is stored under: a hash of what identifies
// the finding, so re-hitting the same bug overwrites its file instead of piling up
// a near-identical copy per page load. The stack is excluded because its line
// numbers shift under edits that do not change the bug, and the commit because it
// changes on every build.
func captureName(c Capture) string {
	h := sha256.New()
	_, _ = fmt.Fprintf(h, "%d\x00%s\x00%s\x00", c.Version, c.PoolDigest, c.Panic)
	// The command log is the finding; JSON of it is stable because every field of
	// input is a scalar or a slice of them.
	if b, err := json.Marshal(c.Snapshot); err == nil {
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))[:16] + ".json"
}

// StoreCapture stamps a posted capture with the dev server's own commit and card
// pool and writes it into dir, returning the path. The server stamps these rather
// than trusting the client because they describe the build the finding must be
// replayed against, which is the server's to know.
func StoreCapture(dir string, body io.Reader) (string, error) {
	var c Capture
	if err := json.NewDecoder(body).Decode(&c); err != nil {
		return "", fmt.Errorf("decode capture: %w", err)
	}
	c.Commit = buildCommit()
	c.PoolDigest = poolDigest()
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, captureName(c))
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// CaptureHandler accepts a posted capture and writes it into dir. It is
// registered only when DevEnabled, so a deployed build has no endpoint at all.
func CaptureHandler(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "post a capture", http.StatusMethodNotAllowed)
			return
		}
		path, err := StoreCapture(dir, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("replay failure captured to %s", path)
		w.WriteHeader(http.StatusNoContent)
	})
}

// postCapture sends a failing match to the dev server so it lands on disk. It is
// fire-and-forget: there is no endpoint outside a dev build and nothing the player
// could do about a failed post, and the quarantine slot still holds the
// reproduction either way.
func postCapture(ctx app.Context, snap snapshot, reason any, stack []byte) {
	if !DevEnabled() {
		return
	}
	body, err := json.Marshal(Capture{
		Version:  snapshotVersion,
		Panic:    fmt.Sprint(reason),
		Stack:    string(stack),
		Snapshot: snap,
	})
	if err != nil {
		return
	}
	ctx.Async(func() {
		resp, err := http.Post(CapturePath, "application/json", bytes.NewReader(body))
		if err != nil {
			return
		}
		_ = resp.Body.Close()
	})
}

// poolDigest hashes the implemented card database: every card's name, type,
// house, power, armor, and rendered text. It is what catches the staleness a
// version bump does not — cards are implemented and reworded constantly, and a
// command log replayed against a pool it was never recorded against deals
// different cards from the same seed, so nothing it does can be reasoned about.
var poolDigest = sync.OnceValue(func() string {
	all := cards.All()
	// The registry hands cards back in no guaranteed order, so sort before hashing
	// or the same pool digests differently from run to run.
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	h := sha256.New()
	for i := range all {
		d := &all[i]
		_, _ = fmt.Fprintf(h, "%s\x00%d\x00%d\x00%d\x00%d\x00%s\x1e",
			d.Name, d.Type, d.House, d.Power, d.Armor, engine.RenderCardText(d))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
})

// buildCommit is the vcs revision this binary was built from, which Go records in
// the build info without any ldflags. It is empty for a binary built outside a
// repository (and for `go test`), which is fine: it is a pointer for whoever
// investigates the finding, not something anything branches on.
func buildCommit() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			return s.Value
		}
	}
	return ""
}
