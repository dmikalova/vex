package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestCaptures replays every capture a dev server has written into
// testdata/capture and fails with the panic it recorded. That is the whole point
// of committing them: a replay failure somebody hit at the keyboard becomes a red
// gate for an agent who was never there. A capture is an open finding, so a green
// tree has none — fix the fault and run `mage capturePrune`.
//
// A capture recorded against a different command-log version or card pool is
// ignored with a note rather than failed. It is not evidence about this tree: the
// same seed deals different cards from a different pool, so nothing the log does
// can be reasoned about. The note lands in the test log rather than on stdout,
// which is where a Go test says something it is not failing over.
func TestCaptures(t *testing.T) {
	caps, err := LoadCaptures(captureTestDir)
	if err != nil {
		t.Fatalf("read the capture directory: %v", err)
	}
	for _, c := range caps {
		t.Run(filepath.Base(c.Path), func(t *testing.T) {
			switch status, detail := c.Replay(); status {
			case CaptureStale:
				t.Skipf("ignoring a stale capture: %s", detail)
			case CaptureFixed:
				t.Logf("no longer reproduces — prune it with `mage capturePrune`")
			case CaptureReproduced:
				t.Errorf("replay failed: %s\nrecorded at commit %q\n%s",
					detail, c.Commit, c.Stack)
			}
		})
	}
}

// failingCapture builds the capture a dev server would have written for a command
// log that panics on replay: a manual move of a card id the deal never handed out,
// the same fault TestAPanickingReplayIsQuarantined quarantines.
func failingCapture(t *testing.T) Capture {
	t.Helper()
	c := newClient(t)
	c.manualTurn(testHouse)
	rec := c.g.s.Record()
	rec.Commands = append(rec.Commands, engine.Command{
		Kind:  engine.CommandManualMove,
		Card:  engine.LocalID(250),
		Index: int(engine.ManualDiscard),
	})
	return Capture{
		Version:    snapshotVersion,
		PoolDigest: poolDigest(),
		Panic:      "index out of range",
		Stack:      "goroutine 1 [running]:\n…",
		Snapshot: snapshot{
			Version: snapshotVersion,
			Record:  rec,
		},
	}
}

// A capture whose command log still fails is reported as reproducing, which is
// what makes TestCaptures red for a finding nobody has fixed yet.
func TestAFailingCaptureReproduces(t *testing.T) {
	status, detail := failingCapture(t).Replay()
	if status != CaptureReproduced {
		t.Fatalf("status = %v, want CaptureReproduced (detail %q)", status, detail)
	}
	if detail == "" {
		t.Error("a reproducing capture reported no detail to debug from")
	}
}

// A capture recorded against another tree is skipped, not failed — for either
// reason it can be stale. Replaying it would prove nothing, and a gate that goes
// red on it would train everybody to delete captures instead of reading them.
func TestAStaleCaptureIsSkipped(t *testing.T) {
	tests := []struct {
		name string
		age  func(c Capture) Capture
	}{
		{"a command log from an older client", func(c Capture) Capture {
			c.Version = snapshotVersion - 1
			return c
		}},
		{"a card pool this tree no longer has", func(c Capture) Capture {
			c.PoolDigest = "0000000000000000"
			return c
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, detail := tt.age(failingCapture(t)).Replay()
			if status != CaptureStale {
				t.Fatalf("status = %v, want CaptureStale", status)
			}
			if detail == "" {
				t.Error("a stale capture reported no reason it was ignored")
			}
		})
	}
}

// A match that replays cleanly is a closed finding: its capture is reported fixed
// so a prune can delete the file.
func TestACaptureThatNoLongerFailsIsFixed(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	rec := Capture{
		Version:    snapshotVersion,
		PoolDigest: poolDigest(),
		Snapshot: snapshot{
			Version: snapshotVersion,
			Record:  c.g.s.Record(),
		},
	}
	if status, detail := rec.Replay(); status != CaptureFixed {
		t.Fatalf("status = %v (detail %q), want CaptureFixed", status, detail)
	}
}

// The round trip a dev server performs: the client posts a failing match, the
// handler stamps it with the build's own commit and card pool and writes it, and
// the replay test reads it back and reproduces the fault. This is the whole
// mechanism end to end, minus the browser.
func TestACapturePostedToTheDevServerIsReplayable(t *testing.T) {
	dir := t.TempDir()
	srv := httptest.NewServer(CaptureHandler(dir))
	defer srv.Close()

	posted := failingCapture(t)
	// The client does not know the server's build, so it posts neither stamp; the
	// handler fills both in.
	posted.PoolDigest = ""
	post(t, srv.URL, posted)

	caps, err := LoadCaptures(dir)
	if err != nil {
		t.Fatalf("load the written capture: %v", err)
	}
	if len(caps) != 1 {
		t.Fatalf("the handler wrote %d captures, want 1", len(caps))
	}
	if caps[0].PoolDigest != poolDigest() {
		t.Errorf("PoolDigest = %q, want the server's %q",
			caps[0].PoolDigest, poolDigest())
	}
	if caps[0].Panic != posted.Panic {
		t.Errorf("Panic = %q, want the posted %q", caps[0].Panic, posted.Panic)
	}
	if status, detail := caps[0].Replay(); status != CaptureReproduced {
		t.Errorf("status = %v (detail %q), want CaptureReproduced", status, detail)
	}
}

// Hitting the same bug twice writes one file, not two: the name is a hash of what
// identifies the finding, so the directory stays a list of open findings rather
// than a tally of how often each was tripped over.
func TestReHittingABugDoesNotPileUpCaptures(t *testing.T) {
	dir := t.TempDir()
	srv := httptest.NewServer(CaptureHandler(dir))
	defer srv.Close()

	c := failingCapture(t)
	post(t, srv.URL, c)
	// A second page load hits the same fault with a stack from a different build.
	c.Stack = "goroutine 7 [running]:\nsomewhere else"
	post(t, srv.URL, c)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read the capture directory: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("the same finding left %d files, want 1", len(entries))
	}
}

// A different finding gets its own file, so deduplicating by content does not
// silently swallow the second bug.
func TestADifferentFindingGetsItsOwnCapture(t *testing.T) {
	dir := t.TempDir()
	srv := httptest.NewServer(CaptureHandler(dir))
	defer srv.Close()

	c := failingCapture(t)
	post(t, srv.URL, c)
	c.Panic = "a different fault entirely"
	post(t, srv.URL, c)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read the capture directory: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("two findings left %d files, want 2", len(entries))
	}
}

// The endpoint takes a posted capture and nothing else: a GET is not a finding,
// and a body that is not a capture is refused rather than written as one.
func TestTheCaptureEndpointRefusesWhatIsNotACapture(t *testing.T) {
	dir := t.TempDir()
	srv := httptest.NewServer(CaptureHandler(dir))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	resp, err = http.Post(srv.URL, "application/json", strings.NewReader("nonsense"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("bad-body status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read the capture directory: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("a refused request wrote %d files, want 0", len(entries))
	}
}

// Pruning keeps the directory a list of open findings: a fault that has been
// fixed and a capture recorded against another tree are both deleted, and only
// the findings that still reproduce are left behind for somebody to fix.
func TestPruningKeepsOnlyOpenFindings(t *testing.T) {
	dir := t.TempDir()
	srv := httptest.NewServer(CaptureHandler(dir))
	defer srv.Close()

	open := failingCapture(t)
	post(t, srv.URL, open)

	// A fixed finding: a command log that replays cleanly now.
	fixed := failingCapture(t)
	fixed.Snapshot.Record.Commands = fixed.Snapshot.Record.Commands[:len(fixed.Snapshot.Record.Commands)-1]
	fixed.Panic = "a fault that has since been fixed"
	post(t, srv.URL, fixed)

	// A stale finding. The handler stamps the current pool, so it is aged on disk.
	stale := failingCapture(t)
	stale.Version = snapshotVersion - 1
	stale.Panic = "a fault from an older client"
	write(t, dir, stale)

	report, err := PruneCaptures(dir)
	if err != nil {
		t.Fatalf("prune: %v", err)
	}
	if report.Scanned != 3 || report.Fixed != 1 || report.Stale != 1 {
		t.Errorf("report = %+v, want 3 scanned, 1 fixed, 1 stale", report)
	}
	if len(report.Open) != 1 {
		t.Fatalf("report left %d open findings, want 1", len(report.Open))
	}

	left, err := LoadCaptures(dir)
	if err != nil {
		t.Fatalf("reload after prune: %v", err)
	}
	if len(left) != 1 {
		t.Fatalf("%d captures survived the prune, want the 1 still failing", len(left))
	}
	if left[0].Panic != open.Panic {
		t.Errorf("the surviving capture is %q, want the still-failing %q",
			left[0].Panic, open.Panic)
	}
}

// write places a capture in dir without going through the handler, which is the
// only way to plant one whose stamps are not this build's.
func write(t *testing.T, dir string, c Capture) {
	t.Helper()
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Fatalf("marshal the capture: %v", err)
	}
	path := filepath.Join(dir, captureName(c))
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatalf("write the capture: %v", err)
	}
}

// post sends one capture to the endpoint the way the client does.
func post(t *testing.T, url string, c Capture) {
	t.Helper()
	body, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal the capture: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post the capture: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("post status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
