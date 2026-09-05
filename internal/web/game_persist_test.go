package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests cover keeping a match alive across a page load: what is written to
// local storage after each action, what is rebuilt from it, and every reason a
// snapshot is thrown away instead.

// A mounted client with nothing saved deals a fresh match rather than coming up
// empty; set selection is reached later through New game.
func TestMountDealsWhenThereIsNothingToResume(t *testing.T) {
	c := newBlankClient(t)
	c.g.OnMount(c.ctx)
	c.settle()

	if c.g.g == nil {
		t.Fatal("mounting dealt no match")
	}
	if c.g.awaitingSetup {
		t.Error("mounting opened the set picker instead of dealing")
	}
	if c.g.phase != phaseHouse {
		t.Errorf("the dealt match is at phase %v, want phaseHouse", c.g.phase)
	}
	if c.g.dispatch == nil {
		t.Error("mounting did not bind the dispatcher the actions complete through")
	}
	if !c.ctx.LocalStorage().Contains(persistKey) {
		t.Error("the dealt match was not saved")
	}
}

// A match saved by one page load comes back on the next, board, log, and view
// together.
func TestAMatchSurvivesAReload(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)
	c.g.sidebarCollapsed = true
	c.g.zonesPlayer = c.g.active()
	c.do(c.g.toggleSidebar)

	before := c.g.g.State
	logLines := len(c.g.g.Log)

	next := c.reload()
	if next.g.seed != c.g.seed {
		t.Errorf("the resumed match has seed %d, want %d", next.g.seed, c.g.seed)
	}
	if next.g.g.State != before {
		t.Error("the resumed match is not the state that was saved")
	}
	if !containsID(next.g.g.Battleline(next.g.active()), id) {
		t.Error("the creature in play did not survive the reload")
	}
	if got := len(next.g.g.Log); got != logLines {
		t.Errorf("the resumed log has %d lines, want %d", got, logLines)
	}
	if next.g.sidebarCollapsed != c.g.sidebarCollapsed {
		t.Error("the sidebar did not come back the way it was left")
	}
	if next.g.zonesPlayer != c.g.active() {
		t.Errorf("the zone viewer came back at %d, want player %d",
			next.g.zonesPlayer, c.g.active())
	}
}

// A PlayerStanding's coloured key tally survives a reload: the typed entry is
// rebuilt from its persisted fields rather than flattened to the plain narrated
// count, so refreshing the page does not change what the key log shows.
func TestAStandingsKeyTallySurvivesAReload(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.g.Log = append(c.g.g.Log, engine.Record{Entry: engine.PlayerStanding{
		Player:    0,
		Aember:    4,
		KeyColors: []engine.KeyColor{engine.KeyColorRed, engine.KeyColorBlue},
	}})
	c.g.save(c.ctx)

	next := c.reload()
	var got engine.PlayerStanding
	var found bool
	for _, rec := range next.g.g.Log {
		if ps, ok := rec.Entry.(engine.PlayerStanding); ok {
			got, found = ps, true
		}
	}
	if !found {
		t.Fatal("the standing did not come back as a typed PlayerStanding")
	}
	if len(got.KeyColors) != 2 ||
		got.KeyColors[0] != engine.KeyColorRed ||
		got.KeyColors[1] != engine.KeyColorBlue {
		t.Errorf("the key tally came back as %v, want two forged keys", got.KeyColors)
	}
}

// A viewer a prompt opened belongs to that prompt, and the prompt does not
// survive the reload, so it comes back closed rather than over nothing.
func TestAPromptsZoneViewerDoesNotSurvive(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.zonesPlayer = c.g.active()
	c.g.promptZone = "discard"
	c.g.save(c.ctx)

	next := c.reload()
	if next.g.zonesPlayer != -1 {
		t.Errorf("the prompt's viewer came back at %d, want closed", next.g.zonesPlayer)
	}
}

// A viewer restored from outside data is range-checked, because a snapshot is
// not the client's own state to trust.
func TestRestoreUIRangeChecksTheViewer(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	for _, in := range []int{-5, 2, 99} {
		c.g.restoreUI(savedUI{ZonesPlayer: in})
		if c.g.zonesPlayer != -1 {
			t.Errorf("a saved viewer of %d restored to %d, want closed",
				in, c.g.zonesPlayer)
		}
	}
}

// The manual card picker only makes sense in manual mode, so it does not come
// back open over an ordinary match.
func TestThePickerOnlyReturnsInManualMode(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.restoreUI(savedUI{PickerOpen: true})
	if c.g.pickerOpen {
		t.Error("the manual picker came back open outside manual mode")
	}
}

// Cards manual mode added are replayed in order, so the rebuilt catalog hands
// out the same ids the saved state refers to.
func TestManualAddsAreReplayed(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.g.save(c.ctx)

	next := c.reload()
	if !containsID(next.g.g.Hand(next.g.active()), id) {
		t.Error("the manually added card did not come back in hand")
	}
	if next.g.g.Def(id).Name != testCreature {
		t.Errorf("id %d rebuilt as %q, want %q", id, next.g.g.Def(id).Name, testCreature)
	}
}

// A snapshot naming a card the pool no longer holds leaves every later id
// misaligned, so it is dropped rather than restored onto the wrong cards.
func TestASnapshotNamingAnUnknownCardIsDropped(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.manualAdds = []manualAdd{{Name: "A Card That Was Never Printed", Player: 0}}
	c.g.save(c.ctx)
	c.expectDropped()
}

// Every reason a snapshot is unusable ends the same way: it is deleted and a
// fresh match is dealt, rather than the client coming up on a state it cannot
// read.
func TestUnusableSnapshotsAreDropped(t *testing.T) {
	tests := []struct {
		name   string
		damage func(snap snapshot) any
	}{
		{"a version from an older engine", func(snap snapshot) any {
			snap.Version = snapshotVersion - 1
			return snap
		}},
		{"no seed to rebuild the catalog from", func(snap snapshot) any {
			snap.Seed = 0
			return snap
		}},
		{"not a snapshot at all", func(_ snapshot) any {
			return "this is not a match"
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t)
			c.startTurn()
			c.g.save(c.ctx)
			var snap snapshot
			if err := c.ctx.LocalStorage().Get(persistKey, &snap); err != nil {
				t.Fatalf("read back the snapshot: %v", err)
			}
			if err := c.ctx.LocalStorage().Set(persistKey, tt.damage(snap)); err != nil {
				t.Fatalf("write the damaged snapshot: %v", err)
			}
			c.expectDropped()
		})
	}
}

// With nothing in storage there is nothing to read, and resume says so without
// touching the slot.
func TestResumeWithNothingSaved(t *testing.T) {
	c := newBlankClient(t)
	if c.g.resume(c.ctx) {
		t.Error("resume restored a match from an empty store")
	}
}

// A half-resolved action cannot be persisted: the rest of it lives on a
// goroutine a reload kills, so the point before it is saved instead and the
// player lands where they can take the action again.
func TestAnActionInFlightSavesThePointBeforeIt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	before := c.g.g.State

	c.g.beginAction()
	c.g.busy = true
	c.g.g.ManualAddAmber(c.g.active(), 5)
	c.g.save(c.ctx)
	c.g.busy = false

	var snap snapshot
	if err := c.ctx.LocalStorage().Get(persistKey, &snap); err != nil {
		t.Fatalf("read back the snapshot: %v", err)
	}
	if snap.State == c.g.g.State {
		t.Error("the half-resolved state was persisted")
	}
	if snap.State != before {
		t.Error("the point before the action in flight was not what was saved")
	}
}

// With no match yet there is nothing to save, and save says so rather than
// writing an empty slot a later load would try to read.
func TestSavingBeforeTheDeal(t *testing.T) {
	c := newBlankClient(t)
	c.g.save(c.ctx)
	if c.ctx.LocalStorage().Contains(persistKey) {
		t.Error("a client with no match wrote a snapshot")
	}
}

// A restored state whose ids do not resolve against the rebuilt catalog would
// panic on the first lookup, so it is caught before the board is drawn.
func TestAStateWithUnreadableIdsIsDropped(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.g.State.Hand[c.g.active()].IDs[0] = engine.LocalID(250)
	c.g.save(c.ctx)
	c.expectDropped()
}

// A resumed match starts its toast caught up: reopening a page with the sidebar
// collapsed must not replay the whole game log into one banner over the board.
func TestReloadingDoesNotToastTheWholeLog(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)      // lines the toast would replay
	c.do(c.g.toggleSidebar) // collapse the sidebar and save it collapsed

	next := c.reload()
	if !next.g.sidebarCollapsed {
		t.Fatal("the reload did not come back with the sidebar collapsed")
	}
	if next.g.toastSeen != len(next.g.g.Log) {
		t.Errorf("toastSeen = %d, want %d (caught up on reload)",
			next.g.toastSeen, len(next.g.g.Log))
	}
	next.g.refreshToast()
	if len(next.g.toastBubbles) != 0 {
		t.Errorf("reloading toasted %d bubbles, want none", len(next.g.toastBubbles))
	}
}

// reload is what a page load does: a fresh component over the same storage,
// resuming what the last one saved.
func (c *client) reload() *client {
	c.t.Helper()
	next := c.nextLoad()
	if !next.g.resume(next.ctx) {
		c.t.Fatal("the saved match was not resumed")
	}
	next.settle()
	return next
}

// expectDropped asserts that the next page load refuses the saved snapshot and
// clears it, so the client deals fresh rather than failing the same way again.
func (c *client) expectDropped() {
	c.t.Helper()
	next := c.nextLoad()
	if next.g.resume(next.ctx) {
		c.t.Error("the unusable snapshot was restored")
	}
	if next.ctx.LocalStorage().Contains(persistKey) {
		c.t.Error("the unusable snapshot was left in storage to fail again")
	}
}
