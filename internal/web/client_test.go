package web

import (
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
	"github.com/dmikalova/vex/internal/session"
)

// This file is the harness the client tests are written against. It drives the
// real root component the way a player does — the same handlers a click or a key
// press calls, over a real engine game — so what these tests pin down is the
// client's actual behaviour rather than a rehearsal of it.
//
// Two facts about go-app shape the harness. First, app.Context is a struct, not
// an interface, so it cannot be faked; one has to be borrowed from the framework.
// Second, go-app only fires OnMount when app.IsClient is true, which it is not in
// a host test, so mounting the component is not what starts it. The way through
// both is a probe component: OnPreRender does run off-browser (app.IsServer is
// true there), and the context it is handed is an ordinary live one, complete
// with a working in-memory local storage and dispatch queue. The harness loads
// the probe, keeps its context, and drives the game component with it.
//
// Everything the client does is SYNCHRONOUS: a click applies a Command to the
// session, which resolves the engine up to the next decision before it returns.
// So settle and await do not wait for anything — they drain go-app's dispatch
// queue and assert what is already true, which is what keeps them worth calling:
// a test that says "await the prompt" fails on the spot if the prompt is not up.
//
// What is out of reach is the DOM: app.Window() reads back empty off-browser, so
// the pieces that measure or scroll elements (the fly-into-play animation, the
// log auto-scroll, the picker's focus) no-op rather than assert. Every one of
// them is written to tolerate a render with no page behind it, which is what
// makes the rest of the client testable here at all.

// testSeed fixes the deal. The decks, the houses, and every card id follow from
// it, so a test can name a card and get the same one every run.
const testSeed = 1

// ctxProbe exists to be handed an app.Context. OnPreRender is the one lifecycle
// hook go-app calls in a host build, so it is the only door a real context comes
// through.
type ctxProbe struct {
	app.Compo
	ctx app.Context
}

func (p *ctxProbe) Render() app.UI { return app.Div() }

func (p *ctxProbe) OnPreRender(ctx app.Context) { p.ctx = ctx }

// client is a mounted game under test: the component, the engine driving its
// dispatch queue, and a context to call handlers with.
type client struct {
	t   *testing.T
	g   *game
	e   app.TestEngine
	ctx app.Context
}

// newClient deals a fresh match from testSeed and leaves it where a player finds
// it: player 0's turn, waiting to choose a house.
func newClient(t *testing.T) *client {
	t.Helper()
	c := newBlankClient(t)
	c.g.dealMatch(testSeed)
	c.keepOpeningHands()
	c.g.inPlayPrev = c.g.inPlaySet()
	return c
}

// keepOpeningHands answers both players' setup mulligan prompts with Keep and
// leaves the match at the house choice — the state a freshly dealt match rests in.
// The prompts arrive one after the other in the turn loop, first player first; a
// test that wants to exercise the mulligan itself drives dealMatch and the prompts
// directly instead of going through newClient.
func (c *client) keepOpeningHands() {
	c.t.Helper()
	for range 2 {
		c.await("a mulligan prompt", c.g.choosingOption)
		c.do(c.g.chooseOptionIdx(0))
	}
	c.await("setup to reach the house choice", func() bool {
		return c.g.phase == phaseHouse
	})
}

// newBlankClient mounts a game with no match dealt, for the tests that want to
// drive the deal itself (OnMount, resume).
func newBlankClient(t *testing.T) *client {
	t.Helper()
	e := app.NewTestEngine()
	p := &ctxProbe{}
	if err := e.Load(p); err != nil {
		t.Fatalf("load probe: %v", err)
	}
	e.ConsumeAll()
	g, ok := NewGame().(*game)
	if !ok {
		t.Fatal("NewGame returned a non-*game")
	}
	c := &client{
		t:   t,
		g:   g,
		e:   e,
		ctx: p.ctx,
	}
	c.g.dispatch = c.ctx.Dispatch
	return c
}

// script stands the client on a session whose driving action is one scripted
// sequence over the board the test has already built, so a test can raise an
// exact prompt and answer it through the real handlers. It is what the old
// background-goroutine chooser rehearsals became: where a match's action is the
// whole turn loop, a test's is a single question — and the client still sees a
// genuine pending Request and answers it with a genuine Command, so what is under
// test is the handoff rather than a rehearsal of it.
//
// The setup closure hands the SAME game back with the board restored to how it
// stood when the script was armed — GameState is a flat value, so putting it back
// is an assignment — which is what lets a rewind replay the script from the same
// starting position the way a rewind of a match replays from the deal.
func (c *client) script(fn func(*engine.Game)) {
	c.t.Helper()
	eg := c.g.eng()
	state := eg.State.FastCopy()
	log := append([]engine.Record(nil), eg.Log...)
	c.g.s = session.New(
		0,
		[2]string{},
		func(int64, [2]string) *engine.Game {
			eg.State = state
			eg.Log = append(eg.Log[:0:0], log...)
			return eg
		},
		fn,
		c.g.lookupCard,
	)
	c.g.rootMarks, c.g.logGroups, c.g.redoLog = nil, nil, nil
	c.g.promptAt = -1
	// The script stands in for the one root action whose resolution raised the
	// prompt, so it is opened like one and can be rewound like one.
	c.g.beginAction()
	c.g.settlePending()
	c.settle()
}

// scriptDone reports whether the scripted action has run to completion, which is
// how a test knows its question was answered and its captured result is final.
func (c *client) scriptDone() bool {
	_, done := c.g.s.Pending()
	return done
}

// settle runs every pending dispatch and deferred operation, so the renders and
// timers a handler queued have run by the time it returns. Nothing about the game
// itself is waited for: applying a Command resolves the engine to its next
// decision before the handler returns.
func (c *client) settle() {
	c.t.Helper()
	c.e.ConsumeAll()
}

// await settles and then asserts cond, which is how a test says what state the
// client should have reached — a prompt raised mid-effect, or the answer to one
// being taken up. It is an assertion, not a wait: the work is already done.
func (c *client) await(what string, cond func() bool) {
	c.t.Helper()
	c.settle()
	if !cond() {
		c.t.Fatalf("the client never reached %s", what)
	}
}

// awaitTimer is await for the handful of things still on a clock: a lift's exit
// animation clearing itself, the selection badges' grow-and-fade. Nothing about
// the game is timed any more — a command resolves before its handler returns —
// so this is only ever about presentation.
func (c *client) awaitTimer(what string, wait time.Duration, cond func() bool) {
	c.t.Helper()
	deadline := time.Now().Add(wait)
	for {
		c.settle()
		if cond() {
			return
		}
		if time.Now().After(deadline) {
			c.t.Fatalf("the client never reached %s", what)
		}
		time.Sleep(time.Millisecond)
	}
}

// do calls a handler the way a click on its element would, then settles.
func (c *client) do(h app.EventHandler) {
	c.t.Helper()
	h(c.ctx, nullEvent())
	c.settle()
}

// nullEvent is a click with no browser behind it. A handler that cancels the
// event or asks what it hit needs a JS value to call into, and off-browser every
// such call reads back empty — which is the same thing a click on nothing means.
func nullEvent() app.Event { return app.Event{Value: app.Null()} }

// nextLoad is the component a page load builds: a new one over the same storage
// and the same dispatch queue, which has not yet been told anything about the
// match the last one was playing.
func (c *client) nextLoad() *client {
	c.t.Helper()
	g, ok := NewGame().(*game)
	if !ok {
		c.t.Fatal("NewGame returned a non-*game")
	}
	next := &client{
		t:   c.t,
		g:   g,
		e:   c.e,
		ctx: c.ctx,
	}
	next.g.dispatch = c.ctx.Dispatch
	return next
}

// press sends a key the way the document listener would, then settles.
func (c *client) press(key string) {
	c.t.Helper()
	c.g.onKey(c.ctx, key, false)
	c.settle()
}

// shiftPress sends a shifted key press.
func (c *client) shiftPress(key string) {
	c.t.Helper()
	c.g.onKey(c.ctx, key, true)
	c.settle()
}

// startTurn answers the opening house prompt with the first house offered, which
// is where every test of ordinary play begins.
func (c *client) startTurn() engine.House {
	c.t.Helper()
	h := c.g.pickableHouses()[0]
	c.do(c.g.pickHouse(h))
	if c.g.phase != phaseMain {
		c.t.Fatalf("after choosing %v the phase is %v, want phaseMain", h, c.g.phase)
	}
	return h
}

// board returns the active player's battleline.
func (c *client) board() []engine.LocalID { return c.g.eng().Battleline(c.g.active()) }

// hand returns the active player's hand.
func (c *client) hand() []engine.LocalID { return c.g.eng().Hand(c.g.active()) }

// manual turns manual mode on, which lifts the house restrictions so a test can
// lay out the board it needs card by card rather than waiting for the deal to
// offer one.
func (c *client) manual() {
	c.t.Helper()
	c.do(c.g.toggleManual)
	if !c.g.eng().Manual() {
		c.t.Fatal("manual mode did not turn on")
	}
}

// manualTurn puts the active player into the main phase of a manual-mode turn
// under house h. It is how a test reaches ordinary play with an exact board
// rather than whatever the deal happened to offer.
func (c *client) manualTurn(h engine.House) {
	c.t.Helper()
	if !c.g.eng().Manual() {
		c.manual()
	}
	c.do(c.g.manualSetHouse(h))
	if c.g.phase != phaseMain {
		c.t.Fatalf("the manual turn is at phase %v, want phaseMain", c.g.phase)
	}
}

// pass ends the turn, clicking through the confirmation a player would see. It
// is how a test gets past summoning sickness: a creature is usable on the turn
// after the one it was played on.
func (c *client) pass() {
	c.t.Helper()
	c.g.confirmEndTurn = true
	c.do(c.g.endTurn)
}

// ownNextTurn hands the turn to the opponent and takes it back, which is how a
// test reaches the turn after the one a creature was played on — the turn it is
// no longer exhausted from entering play.
func (c *client) ownNextTurn(h engine.House) {
	c.t.Helper()
	c.pass()
	c.manualTurn(h)
	c.pass()
	c.manualTurn(h)
}

// deal puts a named card into the active player's hand and returns its id. It
// goes through the client's own manual add, so the command is recorded and a
// reload replays it into the rebuilt catalog.
func (c *client) deal(name string) engine.LocalID {
	c.t.Helper()
	def, ok := c.g.defByName[name]
	if !ok {
		c.t.Fatalf("no card named %q", name)
	}
	player := c.g.active()
	before := len(c.g.eng().Hand(player))
	c.g.applyManual(c.ctx, engine.Command{
		Kind:   engine.CommandManualAddCard,
		Name:   def.Name,
		Player: player,
	})
	hand := c.g.eng().Hand(player)
	if len(hand) != before+1 {
		c.t.Fatalf("%s was not added to hand", name)
	}
	return hand[len(hand)-1]
}

// playFromHand selects a card in hand and plays it, taking the right flank when
// the card is a creature and the battleline is not empty.
func (c *client) playFromHand(id engine.LocalID) {
	c.t.Helper()
	c.g.selectHandID(c.ctx, id)
	if !c.g.hasSel || c.g.sel != id {
		c.t.Fatalf("card %d could not be selected in hand", id)
	}
	c.do(c.g.play)
	if c.g.phase == phaseFlank {
		c.do(c.g.playFlank(false))
	}
	if c.g.status != "" {
		c.t.Fatalf("playing card %d reported %q", id, c.g.status)
	}
}

// promptArtifact is an artifact whose Action ability asks the player to choose a
// creature, so using it raises a mandatory card prompt inside a real root action
// — which is what a test about backing out of a prompt needs to back out of.
const promptArtifact = "Cannon"

// promptHouse is promptArtifact's house, so a manual turn under it can use the
// artifact.
const promptHouse = engine.Brobnar

// stagePromptArtifact puts promptArtifact into play ready to use, and returns its
// id.
func (c *client) stagePromptArtifact() engine.LocalID {
	c.t.Helper()
	id := c.deal(promptArtifact)
	c.g.selectHandID(c.ctx, id)
	c.do(c.g.manualPlay)
	c.g.applyManual(c.ctx, engine.Command{
		Kind: engine.CommandManualReady,
		Card: id,
	})
	c.settle()
	return id
}
