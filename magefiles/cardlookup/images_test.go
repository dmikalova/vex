package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

// testFetcher returns a fetcher for an httptest server: no real backoff.
func testFetcher(dir string) *imageFetcher {
	f := newImageFetcher(dir)
	f.backoff = time.Millisecond
	return f
}

func TestParseImportArgs(t *testing.T) {
	cases := []struct {
		args      []string
		slug, dir string
		wantErr   bool
	}{
		{args: nil},
		{args: []string{"menagerie"}, slug: "menagerie"},
		{args: []string{"all", "-images=tmp/img"}, slug: "all", dir: "tmp/img"},
		{args: []string{"-images=tmp/img", "all"}, slug: "all", dir: "tmp/img"},
		{args: []string{"-images="}, wantErr: true},
		{args: []string{"a", "b"}, wantErr: true},
		{args: []string{"-new"}, wantErr: true},
	}
	for _, c := range cases {
		slug, dir, err := parseImportArgs(c.args)
		if (err != nil) != c.wantErr {
			t.Errorf("parseImportArgs(%q) error = %v, want error %v", c.args, err, c.wantErr)
			continue
		}
		if slug != c.slug || dir != c.dir {
			t.Errorf("parseImportArgs(%q) = %q, %q, want %q, %q", c.args, slug, dir, c.slug, c.dir)
		}
	}
}

func TestImageJobsPaths(t *testing.T) {
	kept := []keptCard{
		{
			Card:     catalogCard{Number: "001", Name: "Alaka"},
			ImageURL: "https://s3/x/341_001_AB_en.PNG",
		},
		{
			Card:     catalogCard{Number: "305", Name: "Veylan Analyst"},
			ImageURL: "https://s3/x/a.png?v=2",
		},
		{Card: catalogCard{Number: "305", Name: "It's Coming..."}, ImageURL: "https://s3/x/b.png"},
		{Card: catalogCard{Number: "S01", Name: "Missing Art"}},
	}
	jobs, noURL, err := imageJobs("out/menagerie", kept)
	if err != nil {
		t.Fatal(err)
	}
	want := []imageJob{
		{URL: "https://s3/x/341_001_AB_en.PNG", Path: filepath.Join("out/menagerie", "001.png")},
		{
			URL:  "https://s3/x/a.png?v=2",
			Path: filepath.Join("out/menagerie", "305-veylan-analyst.png"),
		},
		{URL: "https://s3/x/b.png", Path: filepath.Join("out/menagerie", "305-its-coming.png")},
	}
	if !reflect.DeepEqual(jobs, want) {
		t.Errorf("imageJobs paths =\n%v\nwant\n%v", jobs, want)
	}
	if len(noURL) != 1 || noURL[0].Number != "S01" {
		t.Errorf("imageJobs noURL = %v, want the S01 card", noURL)
	}
}

func TestImageJobsFailsOnPathCollision(t *testing.T) {
	kept := []keptCard{
		{Card: catalogCard{Number: "010", Name: "A-B"}, ImageURL: "https://s3/a.png"},
		{Card: catalogCard{Number: "010", Name: "A B"}, ImageURL: "https://s3/b.png"},
	}
	if _, _, err := imageJobs("out", kept); err == nil {
		t.Error("imageJobs with two cards on one path should error")
	}
}

func TestImageExtNeedsExtension(t *testing.T) {
	if _, err := imageExt("https://s3/media/card"); err == nil {
		t.Error("imageExt with no extension should error")
	}
}

func TestFetchSetDownloadsAndSkipsExisting(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if r.Header.Get("User-Agent") != "vex-cardlookup/1.0" {
			t.Errorf("User-Agent = %q", r.Header.Get("User-Agent"))
		}
		_, _ = io.WriteString(w, "img:"+r.URL.Path)
	}))
	defer srv.Close()

	dir := t.TempDir()
	setDir := filepath.Join(dir, "cota")
	if err := os.MkdirAll(setDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// An image already on disk is kept, and a stale temp file is cleared.
	if err := os.WriteFile(filepath.Join(setDir, "002.png"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(setDir, ".003.png.123.tmp")
	if err := os.WriteFile(stale, []byte("part"), 0o644); err != nil {
		t.Fatal(err)
	}

	kept := []keptCard{
		{Card: catalogCard{Number: "001"}, ImageURL: srv.URL + "/1.png"},
		{Card: catalogCard{Number: "002"}, ImageURL: srv.URL + "/2.png"},
		{Card: catalogCard{Number: "003"}, ImageURL: srv.URL + "/3.png"},
		{Card: catalogCard{Number: "004"}},
	}
	counts, err := testFetcher(dir).fetchSet("cota", kept)
	if err != nil {
		t.Fatal(err)
	}
	if want := (imageCounts{Fetched: 2, Present: 1, NoURL: 1}); counts != want {
		t.Errorf("counts = %+v, want %+v", counts, want)
	}
	if hits.Load() != 2 {
		t.Errorf("server hits = %d, want 2", hits.Load())
	}
	got, err := os.ReadFile(filepath.Join(setDir, "001.png"))
	if err != nil || string(got) != "img:/1.png" {
		t.Errorf("001.png = %q, %v", got, err)
	}
	if got, _ := os.ReadFile(filepath.Join(setDir, "002.png")); string(got) != "old" {
		t.Errorf("002.png was overwritten: %q", got)
	}
	leftovers, _ := filepath.Glob(filepath.Join(setDir, ".*.tmp"))
	if len(leftovers) != 0 {
		t.Errorf("temp files left behind: %v", leftovers)
	}
}

func TestFetchRetries5xxAndFailsOn4xx(t *testing.T) {
	var flaky atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/flaky.png":
			if flaky.Add(1) < 3 {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			_, _ = io.WriteString(w, "ok")
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	dir := t.TempDir()
	f := testFetcher(dir)
	if _, err := f.fetch(
		imageJob{URL: srv.URL + "/flaky.png", Path: filepath.Join(dir, "a.png")},
	); err != nil {
		t.Errorf("flaky 5xx should succeed after retries: %v", err)
	}
	missing := filepath.Join(dir, "b.png")
	if _, err := f.fetch(imageJob{URL: srv.URL + "/gone.png", Path: missing}); err == nil {
		t.Error("a 404 should fail")
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Errorf("a failed download left a file at %s", missing)
	}
}

func TestWriteAtomicLeavesNoFileOnEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()

	dir := t.TempDir()
	f := testFetcher(dir)
	f.retries = 0
	p := filepath.Join(dir, "empty.png")
	if _, err := f.fetch(imageJob{URL: srv.URL + "/e.png", Path: p}); err == nil {
		t.Error("an empty body should fail")
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Errorf("an empty body left a file at %s", p)
	}
	leftovers, _ := filepath.Glob(filepath.Join(dir, ".*.tmp"))
	if len(leftovers) != 0 {
		t.Errorf("temp files left behind: %v", leftovers)
	}
}
