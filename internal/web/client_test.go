package web

import (
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
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
	c.g.inPlayPrev = c.g.inPlaySet()
	return c
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
	c := &client{t: t, g: NewGame().(*game), e: e, ctx: p.ctx}
	c.g.dispatch = c.ctx.Dispatch
	return c
}

// settle runs every pending dispatch, async, and deferred operation, so an action
// that resolves on a background goroutine has finished by the time it returns.
func (c *client) settle() {
	c.t.Helper()
	c.e.ConsumeAll()
}

// await settles until cond holds, which is how a test waits for something a
// background goroutine posted — a prompt raised mid-effect, or the answer to one
// being taken up. It fails rather than spin forever.
func (c *client) await(what string, cond func() bool) {
	c.t.Helper()
	for range 500 {
		c.settle()
		if cond() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	c.t.Fatalf("timed out waiting for %s", what)
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
	next := &client{t: c.t, g: NewGame().(*game), e: c.e, ctx: c.ctx}
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
func (c *client) board() []engine.LocalID { return c.g.g.Battleline(c.g.active()) }

// hand returns the active player's hand.
func (c *client) hand() []engine.LocalID { return c.g.g.Hand(c.g.active()) }

// manual turns manual mode on, which lifts the house restrictions so a test can
// lay out the board it needs card by card rather than waiting for the deal to
// offer one.
func (c *client) manual() {
	c.t.Helper()
	c.do(c.g.toggleManual)
	if !c.g.g.Manual() {
		c.t.Fatal("manual mode did not turn on")
	}
}

// manualTurn puts the active player into the main phase of a manual-mode turn
// under house h. It is how a test reaches ordinary play with an exact board
// rather than whatever the deal happened to offer.
func (c *client) manualTurn(h engine.House) {
	c.t.Helper()
	if !c.g.g.Manual() {
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
// needs manual mode, and records the add the way the card picker does so a
// reload can replay it into the rebuilt catalog.
func (c *client) deal(name string) engine.LocalID {
	c.t.Helper()
	def, ok := c.g.defByName[name]
	if !ok {
		c.t.Fatalf("no card named %q", name)
	}
	player := c.g.active()
	id, added := c.g.g.ManualAddCard(*def, player)
	if !added {
		c.t.Fatalf("%s was not added to hand", name)
	}
	c.g.manualAdds = append(c.g.manualAdds, manualAdd{Name: def.Name, Player: player})
	return id
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
