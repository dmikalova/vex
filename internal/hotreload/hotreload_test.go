package hotreload

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeGo writes a .go file at rel under dir, creating its parent directories,
// and stamps it with modTime so a test can order files without sleeping.
func writeGo(t *testing.T, dir, rel string, modTime time.Time) string {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte("package x\n"), 0o644); err != nil {
		t.Fatalf("WriteFile %s: %v", path, err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("Chtimes %s: %v", path, err)
	}
	return path
}

// TestNewestModTimeSkipsDotDirectories pins the prune that keeps the dev
// server from rebuilding for an edit outside the tree it serves: a newer file
// under a dot-directory, such as an agent worktree in .diatom, does not move
// the reported time.
func TestNewestModTimeSkipsDotDirectories(t *testing.T) {
	root := t.TempDir()
	watched := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	writeGo(t, root, filepath.Join("internal", "engine", "game.go"), watched)

	newer := watched.Add(time.Hour)
	writeGo(t, root, filepath.Join(".diatom", "worktree", "x.go"), newer)
	writeGo(t, root, filepath.Join("node_modules", "pkg", "y.go"), newer)
	writeGo(t, root, filepath.Join("tmp", "z.go"), newer)

	cfg := Config{Root: root, Extensions: []string{".go"}}
	if got := cfg.newestModTime(); !got.Equal(watched) {
		t.Errorf("newestModTime() = %v, want %v (pruned trees must not count)", got, watched)
	}
}

// TestNewestModTimeWalksDotRoot pins the root exemption: filepath.WalkDir
// reports the root entry first, and its name is "." for the default Root, so a
// prune that fired on it would skip the whole tree and look like "nothing
// changed".
func TestNewestModTimeWalksDotRoot(t *testing.T) {
	for _, tc := range []struct {
		name string
		root string
	}{
		{name: "relative dot", root: "."},
		{name: "dot-directory root", root: ".diatom"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			watched := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			writeGo(t, dir, filepath.Join(tc.root, "main.go"), watched)
			t.Chdir(dir)

			cfg := Config{Root: tc.root, Extensions: []string{".go"}}
			if got := cfg.newestModTime(); !got.Equal(watched) {
				t.Errorf("newestModTime() = %v, want %v (the root is always walked)", got, watched)
			}
		})
	}
}

// TestWatches pins which file names trigger a rebuild: a watched extension
// does, and a Go test file never does.
func TestWatches(t *testing.T) {
	cfg := Config{Extensions: []string{".go", ".css"}}
	for name, want := range map[string]bool{
		"game.go":      true,
		"style.css":    true,
		"game_test.go": false,
		"README.md":    false,
	} {
		if got := cfg.watches(name); got != want {
			t.Errorf("watches(%q) = %v, want %v", name, got, want)
		}
	}
}
