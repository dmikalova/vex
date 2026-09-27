//go:build uitest

package uitest

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"

	"github.com/dmikalova/vex/internal/web"
)

// The pace of the run. buildTimeout covers the js/wasm build (~3s warm, longer
// cold); serverTimeout is how long the freshly started server has to answer;
// scenarioTimeout is one whole journey, generous because the first page load
// pulls the 28 MB bundle; statusPoll is how often the status element is read.
const (
	buildTimeout    = 5 * time.Minute
	serverTimeout   = 30 * time.Second
	scenarioTimeout = 90 * time.Second
	statusPoll      = 200 * time.Millisecond
)

// statusJS reads the one element the page writes every outcome to: its
// data-state is "running", "passed", or "failed", and its text names the step
// and why. Reading one place is what keeps the driver free of any knowledge of
// what a scenario does.
const statusJS = `() => {
	const el = document.getElementById('ui-test-status');
	return el ? { state: el.dataset.state, text: el.textContent } : null;
}`

// baseURL is the server the scenarios are driven against and browser the Chrome
// driving them. Both are set up once by TestMain, since a build, a server boot,
// and a browser launch per subtest would dwarf the scenarios themselves.
var (
	baseURL string
	browser *rod.Browser
)

// TestMain builds the client, serves it with the ui-test routes switched on,
// launches headless Chrome, and tears all three down again.
func TestMain(m *testing.M) { os.Exit(run(m)) }

// run is TestMain's body, split out so its deferred teardowns run before the
// process exits.
func run(m *testing.M) int {
	root, err := repoRoot()
	if err != nil {
		return fail("find the repository root", err)
	}
	if err := buildWasm(root); err != nil {
		return fail("build web/app.wasm", err)
	}
	stop, url, err := serve(root)
	if err != nil {
		return fail("serve the client", err)
	}
	defer stop()
	baseURL = url

	l := launcher.New().Headless(true)
	controlURL, err := l.Launch()
	if err != nil {
		return fail("launch headless Chrome", err)
	}
	defer l.Cleanup()
	b := rod.New().ControlURL(controlURL)
	if err := b.Connect(); err != nil {
		return fail("connect to headless Chrome", err)
	}
	defer func() { _ = b.Close() }()
	browser = b

	return m.Run()
}

// fail reports a setup step that did not work and returns the exit code for it.
// A driver that cannot build, serve, or launch has found nothing about the
// client, so it says which step failed rather than failing every scenario.
func fail(what string, err error) int {
	fmt.Fprintf(os.Stderr, "uitest: could not %s: %v\n", what, err)
	return 1
}

// TestScenarios runs one subtest per registered scenario. The list comes from
// internal/web, so this function names no journey and needs no edit when one is
// added.
func TestScenarios(t *testing.T) {
	scenarios := web.UITestScenarios()
	if len(scenarios) == 0 {
		t.Fatal("no browser scenarios are registered")
	}
	for _, s := range scenarios {
		t.Run(s.Slug, func(t *testing.T) {
			state, text := runScenario(t, s.Slug)
			if state != "passed" {
				t.Fatalf("%s: %s", s.Name, text)
			}
		})
	}
}

// runScenario opens one scenario's page and waits for its verdict. The once
// flag stops the page after a single pass, so the driver reads one result
// instead of watching the loop a human watches.
func runScenario(t *testing.T, slug string) (state, text string) {
	t.Helper()
	url := fmt.Sprintf("%s/ui-test/%s?once=1", baseURL, slug)
	page, err := browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		t.Fatalf("open %s: %v", url, err)
	}
	defer func() { _ = page.Close() }()

	deadline := time.Now().Add(scenarioTimeout)
	for {
		state, text = readStatus(page)
		if state == "passed" || state == "failed" {
			return state, text
		}
		if time.Now().After(deadline) {
			if text == "" {
				text = "the page never showed a status element"
			}
			return "timeout", fmt.Sprintf("timed out after %s while %s", scenarioTimeout, text)
		}
		time.Sleep(statusPoll)
	}
}

// readStatus reads the status element, treating a page that has not rendered it
// yet as still running — the bundle is still downloading, or the client has not
// mounted.
func readStatus(page *rod.Page) (state, text string) {
	res, err := page.Eval(statusJS)
	if err != nil || res.Value.Nil() {
		return "running", ""
	}
	return res.Value.Get("state").Str(), res.Value.Get("text").Str()
}

// repoRoot walks up from the test's directory to the one holding go.mod, which
// is the directory cmd/web must be served from (it reads web/app.wasm and
// web/assets by relative path).
func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod above %s", dir)
		}
		dir = parent
	}
}

// buildWasm compiles the client into web/app.wasm, exactly as mage webWasm
// does. The driver builds it rather than assuming a fresh one, so the browser
// runs the working tree and not whatever bundle was last left on disk.
func buildWasm(root string) error {
	cmd := exec.Command(
		"go", "build", "-trimpath", "-ldflags=-s -w", "-o", "web/app.wasm", "./cmd/web")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	return runCmd(cmd, buildTimeout)
}

// serve builds cmd/web natively, starts it on a free port with the ui-test
// routes switched on, and waits for it to answer. It returns the base URL and
// the function that stops it.
func serve(root string) (stop func(), url string, err error) {
	bin := filepath.Join(os.TempDir(), "vex-uitest-web")
	build := exec.Command("go", "build", "-o", bin, "./cmd/web")
	build.Dir = root
	if err := runCmd(build, buildTimeout); err != nil {
		return nil, "", err
	}
	port, err := freePort()
	if err != nil {
		return nil, "", err
	}
	srv := exec.Command(bin)
	srv.Dir = root
	// PORT is the variable cmd/web already reads, so the driver needs no flag of
	// its own; VEX_UITEST is the same switch mage web sets, and the server passes
	// it down to the wasm client so both sides register the routes.
	srv.Env = append(os.Environ(), "PORT="+port, web.UITestEnv+"=1")
	srv.Stdout, srv.Stderr = os.Stderr, os.Stderr
	if err := srv.Start(); err != nil {
		return nil, "", err
	}
	stop = func() {
		_ = srv.Process.Kill()
		_ = srv.Wait()
		_ = os.Remove(bin)
	}
	url = "http://127.0.0.1:" + port
	if err := waitForServer(url); err != nil {
		stop()
		return nil, "", err
	}
	return stop, url, nil
}

// waitForServer polls the scenario index until the server answers it. A 200
// there also proves the routes are switched on: without VEX_UITEST the path is
// not registered and go-app 404s it.
func waitForServer(url string) error {
	deadline := time.Now().Add(serverTimeout)
	var last error
	for time.Now().Before(deadline) {
		resp, err := http.Get(url + "/ui-test")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			err = fmt.Errorf("/ui-test answered %s", resp.Status)
		}
		last = err
		time.Sleep(statusPoll)
	}
	return fmt.Errorf("server did not answer within %s: %w", serverTimeout, last)
}

// freePort asks the kernel for an unused port and hands it back as a string.
func freePort() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer func() { _ = l.Close() }()
	_, port, err := net.SplitHostPort(l.Addr().String())
	return port, err
}

// runCmd runs a build to completion under a time budget, reporting its output
// when it fails so a broken client reads as a compile error rather than a
// browser that found nothing.
func runCmd(cmd *exec.Cmd, timeout time.Duration) error {
	out := make(chan error, 1)
	var buf []byte
	go func() {
		b, err := cmd.CombinedOutput()
		buf = b
		out <- err
	}()
	select {
	case err := <-out:
		if err != nil {
			return fmt.Errorf("%s: %w\n%s", cmd.Args, err, buf)
		}
		return nil
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return fmt.Errorf("%s did not finish within %s", cmd.Args, timeout)
	}
}
