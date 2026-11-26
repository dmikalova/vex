package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// keptCard pairs a catalog card with the front image URL of the feed card it
// came from. The pair is built after mergeCard, so a gigantic creature carries
// the image of its base half.
type keptCard struct {
	Card     catalogCard
	ImageURL string
}

// imageJob is one image to download: its source URL and the file it lands in.
type imageJob struct {
	URL  string
	Path string
}

// imageCounts tallies one set's image run.
type imageCounts struct {
	Fetched int
	Present int
	NoURL   int
}

// imageFetcher downloads card front images. The images are on S3, not the
// throttled decks feed, so a small worker pool fetches them in parallel.
type imageFetcher struct {
	client  *http.Client
	dir     string
	workers int
	retries int
	backoff time.Duration
}

// newImageFetcher returns a fetcher that writes under dir: 4 workers, a 60
// second timeout, and 3 retries that start 2 seconds apart and double.
func newImageFetcher(dir string) *imageFetcher {
	return &imageFetcher{
		client:  &http.Client{Timeout: 60 * time.Second},
		dir:     dir,
		workers: 4,
		retries: 3,
		backoff: 2 * time.Second,
	}
}

// fetchSet downloads the images of one set's kept cards into <dir>/<slug>. It
// first removes temp files an interrupted run left behind, so an upload of the
// directory never carries a partial image. A card with no image URL is logged
// and counted, not fatal.
func (f *imageFetcher) fetchSet(slug string, kept []keptCard) (imageCounts, error) {
	setDir := filepath.Join(f.dir, slug)
	if err := removeStaleTemps(setDir); err != nil {
		return imageCounts{}, err
	}
	jobs, noURL, err := imageJobs(setDir, kept)
	if err != nil {
		return imageCounts{}, err
	}
	for _, c := range noURL {
		fmt.Printf("  no image URL: %s %s\n", c.Number, c.Name)
	}
	counts, err := f.fetchAll(jobs)
	counts.NoURL = len(noURL)
	return counts, err
}

// imageJobs maps a set's kept cards to the files their images land in. A card
// lands at <setDir>/<number>.<ext>, with the extension taken from the image
// URL. The feed can give two different cards one collector number: Menagerie
// lists both Veylan Analyst and It's Coming... as 305. Every card on a shared
// number gets its name slug as a suffix (305-veylan-analyst.png,
// 305-its-coming.png), so neither claims the bare number and the paths do not
// depend on feed order. Two cards that still map to one path are an error. It
// returns the cards with no image URL separately.
func imageJobs(setDir string, kept []keptCard) ([]imageJob, []catalogCard, error) {
	perNumber := map[string]int{}
	for _, k := range kept {
		perNumber[k.Card.Number]++
	}
	var jobs []imageJob
	var noURL []catalogCard
	claimed := map[string]string{}
	for _, k := range kept {
		if k.ImageURL == "" {
			noURL = append(noURL, k.Card)
			continue
		}
		ext, err := imageExt(k.ImageURL)
		if err != nil {
			return nil, nil, fmt.Errorf("card %s %q: %w", k.Card.Number, k.Card.Name, err)
		}
		base := k.Card.Number
		if perNumber[k.Card.Number] > 1 {
			base += "-" + nameSlug(k.Card.Name)
		}
		p := filepath.Join(setDir, base+ext)
		if other, ok := claimed[p]; ok {
			return nil, nil, fmt.Errorf(
				"cards %q and %q both map to image %s",
				other, k.Card.Name, p,
			)
		}
		claimed[p] = k.Card.Name
		jobs = append(jobs, imageJob{URL: k.ImageURL, Path: p})
	}
	return jobs, noURL, nil
}

// imageExt returns the lowercase extension of an image URL's path, with its
// leading dot. A URL whose path has no extension is an error.
func imageExt(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parsing image URL %q: %w", raw, err)
	}
	ext := strings.ToLower(path.Ext(u.Path))
	if ext == "" || ext == "." {
		return "", fmt.Errorf("image URL %q has no file extension", raw)
	}
	return ext, nil
}

// nameSlug turns an ASCII-folded card name into a file name part: lowercase
// letters and digits, with each other run of characters as one hyphen.
// Apostrophes and quotes are dropped, so "It's Coming..." reads "its-coming".
func nameSlug(name string) string {
	var b strings.Builder
	gap := false
	for _, r := range strings.ToLower(name) {
		switch {
		case r == '\'' || r == '"':
			continue // dropped, so "It's" reads "its"
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			if gap && b.Len() > 0 {
				b.WriteByte('-')
			}
			gap = false
			b.WriteRune(r)
		default:
			gap = true
		}
	}
	return b.String()
}

// removeStaleTemps deletes the temp files writeAtomic leaves in dir when a run
// is killed mid-download. A missing dir has nothing to remove.
func removeStaleTemps(dir string) error {
	stale, err := filepath.Glob(filepath.Join(dir, ".*.tmp"))
	if err != nil {
		return err
	}
	for _, p := range stale {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

// fetchAll downloads the jobs with the worker pool. It tries every job, then
// returns the counts and every error joined.
func (f *imageFetcher) fetchAll(jobs []imageJob) (imageCounts, error) {
	queue := make(chan imageJob)
	var (
		mu     sync.Mutex
		counts imageCounts
		errs   []error
		wg     sync.WaitGroup
	)
	for range f.workers {
		wg.Go(func() {
			for j := range queue {
				fetched, err := f.fetch(j)
				mu.Lock()
				switch {
				case err != nil:
					errs = append(errs, err)
				case fetched:
					counts.Fetched++
				default:
					counts.Present++
				}
				mu.Unlock()
			}
		})
	}
	for _, j := range jobs {
		queue <- j
	}
	close(queue)
	wg.Wait()
	return counts, errors.Join(errs...)
}

// fetch downloads one image unless its file already exists with content, so a
// rerun resumes. It retries a network error or a 5xx with doubling backoff, and
// fails at once on any other status. It reports whether it downloaded.
func (f *imageFetcher) fetch(j imageJob) (bool, error) {
	if info, err := os.Stat(j.Path); err == nil && info.Size() > 0 {
		return false, nil
	}
	wait := f.backoff
	for attempt := 0; ; attempt++ {
		retry, err := f.download(j)
		if err == nil {
			return true, nil
		}
		if !retry || attempt >= f.retries {
			return false, err
		}
		time.Sleep(wait)
		wait *= 2
	}
}

// download makes one attempt at an image. It reports whether a failure is
// worth a retry: a network error, a 5xx, or a broken body is; another non-200
// status is not.
func (f *imageFetcher) download(j imageJob) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, j.URL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "vex-cardlookup/1.0")
	resp, err := f.client.Do(req)
	if err != nil {
		return true, fmt.Errorf("fetching %s: %w", j.URL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= http.StatusInternalServerError {
		return true, fmt.Errorf("fetching %s: status %d", j.URL, resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("fetching %s: status %d", j.URL, resp.StatusCode)
	}
	if err := writeAtomic(j.Path, resp.Body); err != nil {
		return true, fmt.Errorf("writing %s: %w", j.Path, err)
	}
	return false, nil
}

// writeAtomic streams r into a temp file next to path and renames it into
// place. An interrupted download leaves only a temp file, never a partial file
// at path. An empty body is an error, because an empty file would look like a
// missing image to the next run anyway.
func writeAtomic(path string, r io.Reader) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	n, err := tmp.ReadFrom(r)
	if err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if n == 0 {
		return errors.New("empty image body")
	}
	return os.Rename(tmp.Name(), path)
}
