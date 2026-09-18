package sim

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// A corpus file is Go's fuzz format, so what SaveCorpus writes has to read back
// as the same bytes the fuzzer would replay.
func TestCorpusFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	script := []byte("\x00\x01\x02\x03\x04\x05\x06\x07quoting \"edge\" \n case")

	path := filepath.Join(dir, "entry")
	body := "go test fuzz v1\n[]byte(" + strconv.Quote(string(script)) + ")\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := readCorpusFile(path)
	if err != nil {
		t.Fatalf("readCorpusFile: %v", err)
	}
	if string(got) != string(script) {
		t.Errorf("round trip = %q, want %q", got, script)
	}
}

// A file with no []byte literal is a corpus entry this package cannot replay, and
// silently treating it as an empty script would quietly drop a regression.
func TestReadCorpusFileRejectsAFileWithNoScript(t *testing.T) {
	path := filepath.Join(t.TempDir(), "entry")
	if err := os.WriteFile(path, []byte("go test fuzz v1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readCorpusFile(path); err == nil {
		t.Error("expected an error for a corpus file with no script")
	}
}

// The numbers in a violation, and any struct it dumps, differ between the scripts
// that trip the same bug, so they must not split one bug across two entries. The
// card name must still tell two bugs apart.
func TestFailureKeyFoldsIncidentalDetail(t *testing.T) {
	a := failureKey(errors.New("creature 48 (Oak) is in play with 5 power and 6 damage"))
	b := failureKey(errors.New("creature 3 (Oak) is in play with 2 power and 4 damage"))
	if a != b {
		t.Errorf("keys differ (%s vs %s); the numbers should be folded away", a, b)
	}
	c := failureKey(errors.New("card 3 (Oak) kept state ({Damage:1 Stunned:false})"))
	d := failureKey(errors.New("card 9 (Oak) kept state ({Damage:2 Stunned:true})"))
	if c != d {
		t.Errorf("keys differ (%s vs %s); a dumped struct should be folded away", c, d)
	}
	if e := failureKey(errors.New("card 3 (Elm) kept state ({Damage:1 Stunned:false})")); c == e {
		t.Error("the same violation on a different card should keep its own entry")
	}
}

// Folding a struct dump must fold its values, not its shape. Two violations that
// dump different fields are different bugs, so they must not collapse onto one
// key and share one corpus entry — that would let the second bug go uncovered.
func TestFailureKeyKeepsDumpedFieldNames(t *testing.T) {
	a := failureKey(errors.New("card 3 (Oak) kept state ({Damage:1 Stunned:false})"))
	b := failureKey(errors.New("card 3 (Oak) kept state ({Damage:1 Warded:false})"))
	if a == b {
		t.Error("dumps naming different fields are different bugs and need different keys")
	}
	c := failureKey(errors.New("card 3 (Oak) kept state ({Damage:1})"))
	if a == c {
		t.Error("dumps of different shapes are different bugs and need different keys")
	}
}

// One bug keeps its shortest reproduction. A later soak finding the same bug with
// a longer script must not bury the minimized entry already on disk.
func TestSaveCorpusKeepsTheShorterReproduction(t *testing.T) {
	dir := t.TempDir()
	// A script that plays out cleanly is left alone by Minimize, so each call files
	// exactly the bytes it is given and the length rule is what is under test.
	short := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	long := []byte{0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}
	violation := errors.New("creature 4 (Oak) is in play with 5 power")

	path, err := SaveCorpus(dir, short, violation)
	if err != nil {
		t.Fatalf("SaveCorpus: %v", err)
	}
	if _, err := SaveCorpus(dir, long, violation); err != nil {
		t.Fatalf("SaveCorpus: %v", err)
	}

	got, err := readCorpusFile(path)
	if err != nil {
		t.Fatalf("readCorpusFile: %v", err)
	}
	if string(got) != string(short) {
		t.Errorf("entry = %q, want the shorter reproduction %q", got, short)
	}
}

// Pruning drops entries whose bug is fixed: a script that plays out cleanly no
// longer pins anything down.
func TestPruneCorpusDropsFixedEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "entry")
	body := "go test fuzz v1\n[]byte(" + strconv.Quote(
		string([]byte{0, 0, 0, 0, 0, 0, 0, 1}),
	) + ")\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := PruneCorpus(dir)
	if err != nil {
		t.Fatalf("PruneCorpus: %v", err)
	}
	if report.Scanned != 1 || report.Fixed != 1 || report.Kept != 0 {
		t.Errorf("report = %+v, want 1 scanned, 1 fixed, 0 kept", report)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("the entry should have been removed")
	}
}
