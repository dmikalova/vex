package web

import (
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// These tests cover the browser-scenario page itself: the registry every
// scenario is declared in, the switch its routes are gated on, the isolation
// that keeps a run off a real match, and what a scenario does when the thing it
// is looking for is not there. The journeys themselves are proved in a browser;
// what is proved here is that the surface running them is sound.

// uiTestURL is a page the host is asked to run, as a browser would hand it one.
func uiTestURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse %q: %v", raw, err)
	}
	return u
}

// newUITestHost mounts a scenario host over the client harness's context, paced
// fast enough that a host run does not sit through the allowances a browser
// needs.
func newUITestHost(t *testing.T, raw string) (*uiTest, *client) {
	t.Helper()
	c := newBlankClient(t)
	u, ok := NewUITest().(*uiTest)
	if !ok {
		t.Fatal("NewUITest returned a non-*uiTest")
	}
	u.timeout = 10 * time.Millisecond
	u.poll = time.Millisecond
	u.pause = time.Millisecond
	u.start(c.ctx, uiTestURL(t, raw))
	c.settle()
	return u, c
}

// Every registered scenario has to be reachable and runnable: a slug the route
// matches, a name to list it under, a seed so each pass deals the same cards, and
// steps that each say what they do. A scenario missing any of those is one the
// page cannot run or a human cannot read.
func TestScenarioRegistryIsWellFormed(t *testing.T) {
	route := regexp.MustCompile(uiTestScenarioPattern)
	seen := map[string]string{}
	for _, s := range uiScenarios {
		if s.Name == "" {
			t.Errorf("scenario %q has no name to list it under", s.Slug)
		}
		if s.Seed == 0 {
			t.Errorf("scenario %q has no seed, so its passes would deal different cards", s.Slug)
		}
		if len(s.Steps) == 0 {
			t.Errorf("scenario %q has no steps", s.Slug)
		}
		if path := uiTestPath + "/" + s.Slug; !route.MatchString(path) {
			t.Errorf("scenario %q is not reachable: %s does not match %s",
				s.Slug, path, uiTestScenarioPattern)
		}
		if prev, dup := seen[s.Slug]; dup {
			t.Errorf("scenarios %q and %q share the slug %q", prev, s.Name, s.Slug)
		}
		seen[s.Slug] = s.Name
		for i, st := range s.Steps {
			if st.Desc == "" {
				t.Errorf("scenario %q step %d has no description", s.Slug, i+1)
			}
			if st.Run == nil {
				t.Errorf("scenario %q step %d does nothing", s.Slug, i+1)
			}
		}
	}
}

// A scenario is looked up by the slug in its path, and an unregistered one is not
// found rather than answered with some other scenario.
func TestScenarioLookupBySlug(t *testing.T) {
	for _, want := range uiScenarios {
		got, ok := uiScenarioBySlug(want.Slug)
		if !ok {
			t.Fatalf("registered scenario %q was not found by its slug", want.Slug)
		}
		if got.Name != want.Name {
			t.Errorf("slug %q resolved to %q, want %q", want.Slug, got.Name, want.Name)
		}
	}
	if _, ok := uiScenarioBySlug("no-such-scenario"); ok {
		t.Error("an unregistered slug resolved to a scenario")
	}
}

// The routes exist only where the switch is set, which is what keeps the
// scenarios out of a production deployment while leaving them compiled into every
// build (they are an environment switch, not a build tag — see uitest.go).
func TestUITestRoutesAreGatedOnTheEnvironment(t *testing.T) {
	if UITestEnabled() {
		t.Fatalf("%s is set in this test process, so the gate cannot be observed", UITestEnv)
	}
	t.Setenv(UITestEnv, "1")
	if !UITestEnabled() {
		t.Errorf("%s=1 did not enable the browser-scenario routes", UITestEnv)
	}
	t.Setenv(UITestEnv, "")
	if UITestEnabled() {
		t.Errorf("%s unset still enabled the browser-scenario routes", UITestEnv)
	}
}

// The per-scenario route matches one path element after /ui-test and nothing
// else, so the index keeps its own route and a deeper path is not served a
// scenario page.
func TestScenarioRoutePatternMatchesOneSlug(t *testing.T) {
	route := regexp.MustCompile(uiTestScenarioPattern)
	for path, want := range map[string]bool{
		"/ui-test/opening":      true,
		"/ui-test/some-slug":    true,
		"/ui-test":              false,
		"/ui-test/":             false,
		"/ui-test/opening/more": false,
		"/play":                 false,
	} {
		if got := route.MatchString(path); got != want {
			t.Errorf("%s matches %s = %v, want %v", uiTestScenarioPattern, path, got, want)
		}
	}
}

// /ui-test lists every registered scenario and links each to its own path, which
// is the whole of what the index is for.
func TestIndexListsEveryScenario(t *testing.T) {
	u, _ := newUITestHost(t, uiTestPath)
	markup := app.HTMLString(u.Render())
	for _, s := range uiScenarios {
		if !strings.Contains(markup, s.Name) {
			t.Errorf("the index does not list %q", s.Name)
		}
		if href := `href="` + uiTestPath + "/" + s.Slug + `"`; !strings.Contains(markup, href) {
			t.Errorf("the index does not link %q (%s)", s.Name, href)
		}
	}
}

// A URL naming a scenario that is not registered fails in the status element,
// where a driver already looks, rather than drawing a blank page it would wait on
// until it timed out.
func TestAnUnknownSlugFailsInTheStatusElement(t *testing.T) {
	u, _ := newUITestHost(t, uiTestPath+"/no-such-scenario?once=1")
	if u.state != uiTestFailed {
		t.Errorf("state = %v, want failed", u.state)
	}
	markup := app.HTMLString(u.Render())
	if !strings.Contains(markup, `id="`+uiTestStatusID+`"`) {
		t.Fatalf("no %s element to read the failure from", uiTestStatusID)
	}
	if !strings.Contains(markup, "no-such-scenario") {
		t.Errorf("the failure does not name the slug that was asked for:\n%s", markup)
	}
}

// A step that cannot find what it is looking for fails saying so. Off-browser
// there is no DOM at all, so every selector comes back empty — which is the
// sharpest form of the case that matters: a scenario run against a board that
// never appeared must report the missing target, never quietly pass.
func TestAStepThatCannotFindItsTargetFails(t *testing.T) {
	u, c := newUITestHost(t, uiTestPath+"/"+uiScenarios[0].Slug+"?once=1")
	deadline := time.Now().Add(5 * time.Second)
	for u.state == uiTestRunning && time.Now().Before(deadline) {
		c.settle()
	}
	if u.state != uiTestFailed {
		t.Fatalf("state = %v, want failed (there is no page for the scenario to click)", u.state)
	}
	if u.passes != 0 {
		t.Errorf("a scenario with no page to click recorded %d passes", u.passes)
	}
	if len(u.results) != 1 || u.results[0].Err == "" {
		t.Fatalf("the failing step was not recorded with its reason: %+v", u.results)
	}
	if !strings.Contains(u.statusText(), u.results[0].Err) {
		t.Errorf("the status does not carry the step's reason: %q", u.statusText())
	}
	// The panel draws the whole run beside the client, so a human watching reads
	// the same failure the driver does.
	markup := app.HTMLString(u.Render())
	for _, want := range []string{
		`id="` + uiTestStatusID + `"`,
		`data-state="failed"`,
		u.scenario.Steps[0].Desc,
		// the reason, which names the selector the step could not find (the markup
		// escapes its quotes, so match the hook itself)
		optionActID("Keep"),
	} {
		if !strings.Contains(markup, want) {
			t.Errorf("the panel does not show %q", want)
		}
	}
	// once=1 means one pass: nothing is scheduled to start another.
	c.settle()
	if u.state != uiTestFailed || u.passes != 0 {
		t.Error("?once=1 started a second pass")
	}
}

// The client a scenario drives saves into its own storage namespace, so a run in
// a real browser cannot overwrite the match a playtester has open — save runs
// after every action, and the first click would otherwise land in the real slot.
func TestTheScenarioClientSavesOutsideTheRealMatch(t *testing.T) {
	c := newBlankClient(t)
	c.g.storeKey = uiTestStoreKey
	c.g.fixedSeed = 7
	c.g.OnMount(c.ctx)
	c.keepOpeningHands()

	if c.g.awaitingSetup {
		t.Error("a seeded client opened the set picker instead of dealing")
	}
	if c.g.seed != 7 {
		t.Errorf("seed = %d, want the injected 7", c.g.seed)
	}
	if !c.ctx.LocalStorage().Contains(uiTestStoreKey) {
		t.Errorf("the match was not saved in %s", uiTestStoreKey)
	}
	if c.ctx.LocalStorage().Contains(persistKey) {
		t.Errorf("the scenario client wrote over the real match in %s", persistKey)
	}
}

// A scenario's selectors are written from the client's data-act hooks, not from
// the words on a button, so this pins the two together: the Keep answer the
// opening scenario clicks is the act value the mulligan's button actually
// carries.
func TestScenarioSelectorsAreWrittenFromTheActHooks(t *testing.T) {
	if got, want := keepSel, `[data-act="option-keep"]`; got != want {
		t.Errorf("the Keep selector is %s, want %s", got, want)
	}
	if got := actSel(optionActID("Keep")); got != keepSel {
		t.Errorf("actSel(optionActID(\"Keep\")) = %s, want %s", got, keepSel)
	}
	c := newBlankClient(t)
	c.g.dealMatch(testSeed)
	c.await("the first mulligan prompt", func() bool { return c.g.choosingOption })
	markup := app.HTMLString(c.g.Render())
	if !strings.Contains(markup, `data-act="option-keep"`) {
		t.Error("the mulligan prompt does not carry the act hook the scenario clicks")
	}
}

// Off-browser every DOM read comes back empty, and the probe reports that as "not
// there" rather than panicking on a page that does not exist — which is what lets
// the host tests above drive the host at all.
func TestThePageProbeReadsAnAbsentPageAsEmpty(t *testing.T) {
	p := &uiPage{}
	if n := p.count(keepSel); n != 0 {
		t.Errorf("count off-browser = %d, want 0", n)
	}
	if _, err := p.find("a Keep button", keepSel); err == nil {
		t.Error("find reported an element on a page that does not exist")
	}
	if err := p.click("a Keep button", keepSel); err == nil {
		t.Error("click reported a press on a page that does not exist")
	}
	if err := p.absent("a Keep button", keepSel); err != nil {
		t.Errorf("absent on an empty page = %v, want nil", err)
	}
}
