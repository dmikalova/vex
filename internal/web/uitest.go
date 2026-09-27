package web

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the browser-scenario page: /ui-test lists the scenarios
// registered in uitest_scenarios.go, and /ui-test/<slug> runs one of them
// against the live client in a real browser, beside a per-step pass/fail panel.
// The scenarios themselves are data; this is only the host that drives them.
//
// The page lives in package web rather than a child package because a scenario
// clicks the client's own DOM: the ids and data-act hooks it selects by are the
// unexported boardCardID/handCardID/act* names, and the host reaches into the
// unexported game component to deal its match. A child package that knew what to
// click would have to be imported back by the very package it drives.
//
// It is not behind a build tag, for the reason recorded in
// docs/adr/0014-style-gallery-on-real-components.md: go-app routes on the
// client, so a tag would exclude the scenarios from the wasm bundle every player
// downloads — which means they would not be compiled by default and would rot as
// the //go:build todo card stubs do. The routes are an environment switch
// instead, off in production, registered on both sides of the build.

// UITestEnv is the variable that turns the browser-scenario page on. mage web
// sets it, so the scenarios exist wherever the client is being developed and
// nowhere else.
const UITestEnv = "VEX_UITEST"

// UITestEnabled reports whether the browser-scenario routes should exist. Like
// the style gallery's switch it has to agree on both sides of the build — the
// server must register the route or it serves no page at all, and the wasm client
// must register it or the served page renders nothing — which app.Getenv bridges
// by reading the process environment on the server and the Env map the server
// passed down on the client.
func UITestEnabled() bool { return app.Getenv(UITestEnv) == "1" }

// uiTestPath is the index listing every scenario, and uiTestScenarioPattern the
// per-scenario page matching exactly one path element after it. The slug is in
// the path, not the query, so the hot reload that follows an edit comes back to
// the same looping scenario.
const (
	uiTestPath            = "/ui-test"
	uiTestScenarioPattern = `^/ui-test/[^/]+$`
)

// UITestRoutes registers the browser-scenario routes when the switch is on. Both
// cmd/web's server build and its wasm build call it, so one function decides the
// routes on both sides.
func UITestRoutes() {
	if !UITestEnabled() {
		return
	}
	app.Route(uiTestPath, NewUITest)
	app.RouteWithRegexp(uiTestScenarioPattern, NewUITest)
}

// uiTestStoreKey is the scratch local-storage slot every scenario's client saves
// into. game_persist's save runs after every action, so without a namespace of
// its own a scenario run in a real browser would overwrite the playtester's open
// match on its first click. One slot is shared by all scenarios: what is being
// isolated is a real game from the scenarios, not the scenarios from each other.
const uiTestStoreKey = "vex.uitest"

// uiTestStatusID names the one element a driver reads: its data-state is
// "running", "passed", or "failed", and its text says which step and why. Every
// outcome lands there, so a driver polls one place rather than scraping the
// panel.
const uiTestStatusID = "ui-test-status"

// The default pace of a run. stepTimeout is how long a step may keep failing
// before the pass is failed — long enough for a deal, an animation, or an engine
// action resolving on a background goroutine, short enough that a genuinely
// missing target is reported rather than waited on. stepPoll is how often a step
// is retried, and loopPause how long a finished pass rests before the next one
// starts.
const (
	stepTimeout = 3 * time.Second
	stepPoll    = 100 * time.Millisecond
	loopPause   = 2 * time.Second
)

// uiTestState is how the current pass stands.
type uiTestState uint8

const (
	uiTestRunning uiTestState = iota
	uiTestPassed
	uiTestFailed
)

// String is the value a driver matches on, and the modifier the panel is styled
// by.
func (s uiTestState) String() string {
	switch s {
	case uiTestPassed:
		return "passed"
	case uiTestFailed:
		return "failed"
	case uiTestRunning:
		return "running"
	}
	return "running"
}

// uiStepResult is one finished step in the panel: what it was for, and the reason
// it failed ("" when it passed).
type uiStepResult struct {
	Desc string
	Err  string
}

// NewUITest returns the root component for the browser-scenario page. The
// scenario it runs comes from the URL at mount, so both routes share it.
func NewUITest() app.Composer {
	return &uiTest{
		timeout: stepTimeout,
		poll:    stepPoll,
		pause:   loopPause,
	}
}

// uiTest is the scenario host: it owns the client the scenario drives and the
// run it is part way through.
type uiTest struct {
	app.Compo

	// slug is the scenario asked for by the URL ("" on the index page), scenario
	// the one it named, and found whether there was one.
	slug     string
	scenario uiScenario
	found    bool
	// once stops after a single pass (?once=1), which is how a driver reads one
	// result instead of watching a loop.
	once bool

	// client is the live client the scenario clicks, dealt from the scenario's seed
	// into the ui-test storage namespace.
	client *game

	// step is the step being attempted, results the ones already finished, and
	// deadline when the current step stops being retried.
	step     int
	results  []uiStepResult
	deadline time.Time
	// state is how the pass stands, failure why it failed, and passes how many have
	// completed — the loop counter the panel shows so a stuck run is visible.
	state   uiTestState
	failure string
	passes  int
	// gen identifies the current pass, so the scheduled retry of an abandoned one
	// (a new pass, or a dismount) returns instead of driving the run that replaced
	// it.
	gen int

	// timeout, poll, and pause are the run's pace, defaulted by NewUITest and
	// shortened by the tests so a host run does not sit through a browser's
	// allowances.
	timeout, poll, pause time.Duration
}

// OnMount starts the run the URL asked for.
func (u *uiTest) OnMount(ctx app.Context) { u.start(ctx, ctx.Page().URL()) }

// OnDismount abandons the run, so a scheduled retry does not keep clicking a page
// that has been navigated away from.
func (u *uiTest) OnDismount() { u.gen++ }

// start reads the scenario and the once flag out of the page URL and, when the
// slug names a registered scenario, builds the client and begins the first pass.
// It takes the URL rather than reading it so a test can drive a page the host
// build has no browser to navigate.
func (u *uiTest) start(ctx app.Context, pageURL *url.URL) {
	u.slug = strings.Trim(strings.TrimPrefix(pageURL.Path, uiTestPath), "/")
	u.once = pageURL.Query().Get("once") == "1"
	if u.slug == "" {
		return
	}
	u.scenario, u.found = uiScenarioBySlug(u.slug)
	if !u.found {
		u.state = uiTestFailed
		u.failure = fmt.Sprintf("no scenario has the slug %q", u.slug)
		return
	}
	u.client = u.newClient()
	u.beginPass(ctx)
}

// newClient builds the client a scenario drives. It deals itself from the
// scenario's seed on its own mount (it has an injected seed, so it skips the set
// picker), into the scratch storage slot so the run cannot touch a real match —
// and, on a second one built mid-run, resumes the match the first one saved
// there, which is how the reload journey gets its page load.
func (u *uiTest) newClient() *game {
	client := newGame()
	client.storeKey = uiTestStoreKey
	client.fixedSeed = u.scenario.Seed
	return client
}

// beginPass resets the board and starts the scenario at its first step. The reset
// is what makes a scenario repeatable: one that plays a creature cannot run again
// on the board it left behind, so each pass clears the ui-test storage slot and
// re-deals from the scenario's seed.
func (u *uiTest) beginPass(ctx app.Context) {
	u.gen++
	u.results = nil
	u.step = 0
	u.state = uiTestRunning
	u.failure = ""
	ctx.LocalStorage().Del(uiTestStoreKey)
	// A scenario that closed the match (the reload journey) may have left the pass
	// with no client at all, so the next pass builds one. It mounts into a cleared
	// storage slot, so it deals the scenario's seed rather than resuming.
	if u.client == nil {
		u.client = u.newClient()
	}
	if u.client.dispatch != nil {
		u.client.dealMatch(u.scenario.Seed)
	}
	// On the first pass the client has not mounted yet — it mounts on the render
	// this one schedules — and its own OnMount deals the seeded match, so there is
	// nothing to re-deal and the first step simply waits for the board.
	u.armStep(ctx, u.gen)
}

// armStep gives the current step its retry budget and schedules its first
// attempt.
func (u *uiTest) armStep(ctx app.Context, gen int) {
	u.deadline = time.Now().Add(u.timeout)
	ctx.After(u.poll, func(c app.Context) { u.attempt(c, gen) })
}

// attempt runs the current step once. A step that does not hold yet is retried
// until its budget runs out, because a click lands before the render, the engine
// action, and the animation it sets off have finished; only a step still failing
// at the deadline fails the pass, and it fails with the reason the step itself
// gave.
func (u *uiTest) attempt(ctx app.Context, gen int) {
	if gen != u.gen {
		return
	}
	if u.step >= len(u.scenario.Steps) {
		u.finishPass(ctx, gen)
		return
	}
	st := u.scenario.Steps[u.step]
	err := st.Run(&uiPage{host: u})
	switch {
	case err == nil:
		u.results = append(u.results, uiStepResult{Desc: st.Desc})
		u.step++
		u.armStep(ctx, gen)
	case time.Now().Before(u.deadline):
		ctx.After(u.poll, func(c app.Context) { u.attempt(c, gen) })
	default:
		u.results = append(u.results, uiStepResult{Desc: st.Desc, Err: err.Error()})
		u.state = uiTestFailed
		u.failure = err.Error()
		u.loop(ctx, gen)
	}
}

// finishPass records a completed pass and loops.
func (u *uiTest) finishPass(ctx app.Context, gen int) {
	u.passes++
	u.state = uiTestPassed
	u.loop(ctx, gen)
}

// loop waits out the pause and runs the scenario again, unless a single pass was
// asked for. A failed pass loops too: the page is what a human watches while they
// fix the client, and the run going green is how they see the fix land.
func (u *uiTest) loop(ctx app.Context, gen int) {
	if u.once {
		return
	}
	ctx.After(u.pause, func(c app.Context) {
		if gen != u.gen {
			return
		}
		u.beginPass(c)
	})
}

// Render draws the index when no scenario was named, and the scenario beside its
// panel when one was.
func (u *uiTest) Render() app.UI {
	if !u.found {
		return u.indexView()
	}
	return app.Div().Class("ui-test").Body(
		app.Div().Class("ui-test-panel").Body(u.panelView()),
		// The client is drawn conditionally so dropping it really takes it out of
		// the tree: go-app patches a component of the same type in place, so a
		// replaced *game would keep the old one's mounted state.
		app.Div().Class("ui-test-client").Body(
			app.If(u.client != nil, func() app.UI { return u.client }),
		),
	)
}

// indexView lists the registered scenarios. A slug that named none is reported
// here too, through the same status element, so a driver pointed at a stale URL
// reads a failure rather than a blank page.
func (u *uiTest) indexView() app.UI {
	return app.Div().Class("ui-test", "ui-test--index").Body(
		app.Div().Class("ui-test-panel").Body(
			app.H1().Class("ui-test-title").Text("Browser scenarios"),
			app.If(u.slug != "", func() app.UI { return u.statusView() }),
			app.Ul().Class("ui-test-list").Body(
				app.Range(uiScenarios).Slice(func(i int) app.UI {
					s := uiScenarios[i]
					return app.Li().Body(
						app.A().Href(uiTestPath + "/" + s.Slug).Text(s.Name),
					)
				}),
			),
		),
	)
}

// panelView is the running scenario's report: its name, the status element, and
// one row per step.
func (u *uiTest) panelView() app.UI {
	return app.Div().Body(
		app.H1().Class("ui-test-title").Text(u.scenario.Name),
		app.A().Class("ui-test-back").Href(uiTestPath).Text("all scenarios"),
		u.statusView(),
		app.Ol().Class("ui-test-steps").Body(
			app.Range(u.scenario.Steps).Slice(func(i int) app.UI {
				return u.stepView(i)
			}),
		),
	)
}

// stepView draws one step: what it does, how it went, and — when it failed — what
// it could not find.
func (u *uiTest) stepView(i int) app.UI {
	mark, cls, why := "·", "ui-test-step--pending", ""
	if i < len(u.results) {
		mark, cls = "✓", "ui-test-step--passed"
		if why = u.results[i].Err; why != "" {
			mark, cls = "✗", "ui-test-step--failed"
		}
	} else if i == u.step && u.state == uiTestRunning {
		cls = "ui-test-step--running"
	}
	return app.Li().Class("ui-test-step", cls).Body(
		app.Span().Class("ui-test-mark").Text(mark),
		app.Span().Class("ui-test-desc").Text(u.scenario.Steps[i].Desc),
		app.If(why != "", func() app.UI {
			return app.Span().Class("ui-test-why").Text(why)
		}),
	)
}

// statusView is the one element a driver reads.
func (u *uiTest) statusView() app.UI {
	return app.Div().
		ID(uiTestStatusID).
		Class("ui-test-status", "ui-test-status--"+u.state.String()).
		DataSet("state", u.state.String()).
		DataSet("passes", strconv.Itoa(u.passes)).
		Text(u.statusText())
}

// statusText says where the run stands in one sentence, so the element a driver
// reads also tells a human watching the page what happened.
func (u *uiTest) statusText() string {
	steps := len(u.scenario.Steps)
	switch u.state {
	case uiTestPassed:
		return fmt.Sprintf("passed all %d steps (pass %d)", steps, u.passes)
	case uiTestFailed:
		if !u.found {
			return u.failure
		}
		return fmt.Sprintf("failed at step %d of %d: %s — %s",
			u.step+1, steps, u.scenario.Steps[u.step].Desc, u.failure)
	case uiTestRunning:
		if u.step >= steps {
			return "running"
		}
		return fmt.Sprintf("running step %d of %d: %s",
			u.step+1, steps, u.scenario.Steps[u.step].Desc)
	}
	return "running"
}

// uiPage is a scenario step's view of the live page: the few DOM reads and the
// one click every step is written from. Selectors are written from the client's
// own hooks — the data-act values and the card element ids — so a scenario keeps
// working through a relabel or a restyle, which a match on rendered text would
// not.
//
// Off-browser every read comes back empty, which is what a host test sees: a step
// then fails naming the target it could not find, which is the right answer for
// "there is no page here" and the reason a scenario can never silently pass.
//
// host is the scenario host behind the page, which the two page-lifetime steps
// (dropClient/mountClient) act on: tearing the client down and standing a fresh
// one up is the one thing no click can do, because a real browser reload would
// restart the run rather than the match. It is nil in a probe built by hand,
// which those two steps report rather than panic on.
type uiPage struct{ host *uiTest }

// actSel matches the control carrying a data-act hook, which is how a scenario
// says what it is clicking rather than which words are on it.
func actSel(act string) string { return `[data-act="` + act + `"]` }

// doc is the page the selectors run against.
func (p *uiPage) doc() app.Value { return app.Window().Get("document") }

// find returns the first element matching sel, or an error naming what was being
// looked for. what is written as the noun phrase a player would use ("a playable
// creature in hand"), because it is what the panel and the status element show.
func (p *uiPage) find(what, sel string) (app.Value, error) {
	el := p.doc().Call("querySelector", sel)
	if !el.Truthy() {
		return el, fmt.Errorf("%s is not on the page (%s matched nothing)", what, sel)
	}
	return el, nil
}

// count reports how many elements match sel.
func (p *uiPage) count(sel string) int {
	return p.doc().Call("querySelectorAll", sel).Get("length").Int()
}

// click presses the first element matching sel. A disabled control is a failure
// rather than a silent no-op: the click would do nothing and the step after it
// would fail somewhere less informative.
func (p *uiPage) click(what, sel string) error {
	el, err := p.find(what, sel)
	if err != nil {
		return err
	}
	if el.Get("disabled").Truthy() {
		return fmt.Errorf("%s is disabled (%s)", what, sel)
	}
	el.Call("click")
	return nil
}

// fill types text into a text box: it sets the value and fires the input event
// the client listens for, which is what a real keystroke does. It is how a
// scenario searches the manual card picker for the card it means to add.
func (p *uiPage) fill(what, sel, text string) error {
	el, err := p.find(what, sel)
	if err != nil {
		return err
	}
	el.Set("value", text)
	el.Call("dispatchEvent", app.Window().Get("Event").New("input"))
	return nil
}

// dropClient tears the match down the way closing the page does: the client
// component leaves the tree, so its state goes with it and only what it wrote to
// local storage is left. It is idempotent, because a step is retried until its
// check holds.
func (p *uiPage) dropClient() error {
	if p.host == nil {
		return fmt.Errorf("there is no scenario host to close the match in")
	}
	p.host.client = nil
	return nil
}

// mountClient stands a fresh client up over the same storage slot, the way
// loading the page again does: a new component resumes the saved match from the
// command log rather than being handed the old one's state. It is idempotent, so
// the step that checks the match came back does not build a second client.
func (p *uiPage) mountClient() error {
	if p.host == nil {
		return fmt.Errorf("there is no scenario host to reopen the match in")
	}
	if p.host.client == nil {
		p.host.client = p.host.newClient()
	}
	return nil
}

// absent is the opposite check: that something a step took away is gone.
func (p *uiPage) absent(what, sel string) error {
	if n := p.count(sel); n > 0 {
		return fmt.Errorf("%s is still on the page (%s matched %d)", what, sel, n)
	}
	return nil
}
