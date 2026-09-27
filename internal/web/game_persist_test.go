package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
	"github.com/dmikalova/vex/internal/session"
)

// These tests cover keeping a match alive across a page load: what is written to
// local storage after each action, what is rebuilt from it, and every reason a
// snapshot is thrown away instead.

// A mounted client with nothing saved opens the new-game set picker rather than
// silently dealing the base set: a first-time visit chooses its two sets, and
// only then is the first match dealt. (Previously mounting dealt a base-set game;
// that assumed CotA on a fresh load, which is the behavior this now replaces.)
func TestMountOpensSetPickerWhenThereIsNothingToResume(t *testing.T) {
	c := newBlankClient(t)
	c.g.OnMount(c.ctx)
	c.settle()

	if !c.g.awaitingSetup {
		t.Error("mounting did not open the set picker")
	}
	if c.g.s != nil {
		t.Error("mounting dealt a match instead of waiting for a set choice")
	}
	if c.g.dispatch == nil {
		t.Error("mounting did not bind the dispatcher the actions complete through")
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

	before := c.g.eng().State
	logLines := len(c.g.eng().Log)

	next := c.reload()
	if next.g.seed != c.g.seed {
		t.Errorf("the resumed match has seed %d, want %d", next.g.seed, c.g.seed)
	}
	if next.g.eng().State != before {
		t.Error("the resumed match is not the state that was saved")
	}
	if !containsID(next.g.eng().Battleline(next.g.active()), id) {
		t.Error("the creature in play did not survive the reload")
	}
	if got := len(next.g.eng().Log); got != logLines {
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

// A PlayerStanding's coloured key tally survives a reload: replaying the command
// log regenerates the real typed end-of-turn standing rather than a flattened
// count, so refreshing the page does not change what the key log shows.
func TestAStandingsKeyTallySurvivesAReload(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	// Forge a red key for the active player, then end the turn so the end-of-turn
	// standing records it. Replay on reload rebuilds that standing, colours and all.
	me := c.g.active()
	c.g.forgingKey = me
	c.do(c.g.pickForgeColor(engine.KeyColorRed))
	c.pass()

	next := c.reload()
	var got engine.PlayerStanding
	var found bool
	for _, rec := range next.g.eng().Log {
		if ps, ok := rec.Entry.(engine.PlayerStanding); ok && ps.Player == me {
			got, found = ps, true
		}
	}
	if !found {
		t.Fatal("the standing did not come back as a typed PlayerStanding")
	}
	if len(got.KeyColors) != 1 || got.KeyColors[0] != engine.KeyColorRed {
		t.Errorf("the key tally came back as %v, want one red key", got.KeyColors)
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
	if !containsID(next.g.eng().Hand(next.g.active()), id) {
		t.Error("the manually added card did not come back in hand")
	}
	if next.g.eng().Def(id).Name != testCreature {
		t.Errorf("id %d rebuilt as %q, want %q", id, next.g.eng().Def(id).Name, testCreature)
	}
}

// A command log naming a card the pool no longer holds cannot be replayed — the
// missing registration leaves every later id misaligned — so the match is not
// restored onto the wrong cards. The snapshot is current-version, so it is
// quarantined rather than deleted.
func TestASnapshotNamingAnUnknownCardIsQuarantined(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.saveWithExtraCommand(engine.Command{
		Kind:   engine.CommandManualAddCard,
		Name:   "A Card That Was Never Printed",
		Player: 0,
	})
	c.expectQuarantined()
}

// saveWithExtraCommand saves the match and then splices cmd onto the end of the
// saved command log, which is how a test stages a record this build cannot
// replay: the live session would refuse the command outright, so it is added to
// the persisted record rather than applied.
func (c *client) saveWithExtraCommand(cmd engine.Command) {
	c.t.Helper()
	c.g.save(c.ctx)
	var snap snapshot
	if err := c.ctx.LocalStorage().Get(c.g.matchKey(), &snap); err != nil {
		c.t.Fatalf("read back the snapshot: %v", err)
	}
	c.writeSnapshotWithCommand(snap, cmd)
}

// writeSnapshotWithCommand splices cmd onto base's command log and writes the
// result to the client's match key, returning the snapshot written. A caller
// stages a second failing record on top of a first this way — rather than
// re-reading storage, which a failed resume has already cleared — so building up
// a growing record across several failures does not depend on the live slot
// still holding the previous one.
func (c *client) writeSnapshotWithCommand(base snapshot, cmd engine.Command) snapshot {
	c.t.Helper()
	base.Record.Commands = append(base.Record.Commands, cmd)
	if err := c.ctx.LocalStorage().Set(c.g.matchKey(), base); err != nil {
		c.t.Fatalf("write the spliced snapshot: %v", err)
	}
	return base
}

// A snapshot rejected before any replay is attempted — a version this build no
// longer reads, a seed it cannot re-deal from, a payload that is not a snapshot —
// is deleted outright and a fresh match is dealt. These are the innocent
// explanations, so none of them is evidence and none is quarantined; only a
// snapshot that got past these checks and then failed to replay is kept (see
// TestAPanickingReplayIsQuarantined).
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
			snap.Record.Seed = 0
			return snap
		}},
		{"a command log from another session version", func(snap snapshot) any {
			snap.Record.Version = session.Version + 1
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
			next := c.expectDropped()
			if next.ctx.LocalStorage().Contains(quarantineKey) {
				t.Error("a snapshot rejected before replay was quarantined, not deleted")
			}
			if next.g.notice == replayFailedNotice {
				t.Error("a snapshot rejected before replay raised the replay-failed notice")
			}
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

// An action stopped at a prompt is saved in full, prompt and all. Nothing is held
// back any more: a prompt's answers so far are ordinary recorded commands, so the
// reload replays them and lands the player back on the same question.
func TestAMatchStoppedAtAPromptResumesAtThatPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	id := c.deal(deployCreature)
	c.g.selectHandID(c.ctx, id)
	c.do(c.g.play)
	c.await("the deploy placement prompt", c.g.choosingPosition)

	saved := len(c.g.s.Record().Commands)
	next := c.reload()
	if got := len(next.g.s.Record().Commands); got != saved {
		t.Errorf("the resumed log holds %d commands, want the %d that were saved",
			got, saved)
	}
	if !next.g.choosingPosition() {
		t.Error("the reload did not come back on the placement prompt")
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

// A command log that replays into a corrupt engine state — here a manual move of
// a card id the deal never handed out — panics during replay. That panic is
// caught and the board is not drawn from a broken state, but the snapshot is
// kept: resume has already ruled out a wrong version and a bad decode, so what is
// left is a save this build cannot replay, and deleting it would destroy the
// reproduction. Only the latest one is kept — the slot is overwritten, not
// appended to.
func TestAPanickingReplayIsQuarantined(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.saveWithExtraCommand(engine.Command{
		Kind:  engine.CommandManualMove,
		Card:  engine.LocalID(250),
		Index: int(engine.ManualDiscard),
	})
	first := c.expectQuarantined()
	if n := len(first.Record.Commands); n != len(c.g.s.Record().Commands)+1 {
		t.Errorf("quarantined %d commands, want the failing log's %d",
			n, len(c.g.s.Record().Commands)+1)
	}

	// A second failure replaces the first rather than piling up beside it.
	c.writeSnapshotWithCommand(first, engine.Command{
		Kind:  engine.CommandManualMove,
		Card:  engine.LocalID(251),
		Index: int(engine.ManualDiscard),
	})
	second := c.expectQuarantined()
	if len(second.Record.Commands) != len(first.Record.Commands)+1 {
		t.Errorf("quarantined %d commands, want the later log's %d",
			len(second.Record.Commands), len(first.Record.Commands)+1)
	}
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
	if next.g.toastSeen != len(next.g.eng().Log) {
		t.Errorf("toastSeen = %d, want %d (caught up on reload)",
			next.g.toastSeen, len(next.g.eng().Log))
	}
	next.g.refreshToast()
	if len(next.g.toastBubbles) != 0 {
		t.Errorf("reloading toasted %d bubbles, want none", len(next.g.toastBubbles))
	}
}

// Dismissing the toast clears every bubble and catches the seen floor up, so the
// same lines do not surface again on the next refresh.
func TestDismissingTheToastClearsIt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.do(c.g.toggleSidebar) // collapse the sidebar so lines toast
	id := c.deal(testCreature)
	c.playFromHand(id)
	c.g.refreshToast()
	if len(c.g.toastBubbles) == 0 {
		t.Fatal("expected a toast bubble to surface")
	}

	c.g.clearToast()
	if len(c.g.toastBubbles) != 0 {
		t.Errorf("dismiss left %d bubbles, want 0", len(c.g.toastBubbles))
	}
	if c.g.toastSeen != len(c.g.eng().Log) {
		t.Errorf("toastSeen = %d, want %d (caught up on dismiss)",
			c.g.toastSeen, len(c.g.eng().Log))
	}
	c.g.refreshToast()
	if len(c.g.toastBubbles) != 0 {
		t.Errorf("dismissed lines toasted again: %d bubbles", len(c.g.toastBubbles))
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
// clears the live slot, so the client deals fresh rather than failing the same way
// again. It hands that page load back so a caller can go on to assert where the
// snapshot went.
func (c *client) expectDropped() *client {
	c.t.Helper()
	next := c.nextLoad()
	if next.g.resume(next.ctx) {
		c.t.Error("the unusable snapshot was restored")
	}
	if next.ctx.LocalStorage().Contains(persistKey) {
		c.t.Error("the unusable snapshot was left in storage to fail again")
	}
	return next
}

// expectQuarantined asserts that the next page load refuses the saved snapshot,
// keeps it under the quarantine key instead of deleting it, and tells the player
// the save was kept. It returns the snapshot that was kept.
func (c *client) expectQuarantined() snapshot {
	c.t.Helper()
	next := c.expectDropped()
	if next.g.notice != replayFailedNotice {
		c.t.Errorf("notice = %q, want the replay-failed notice", next.g.notice)
	}
	var kept snapshot
	if err := next.ctx.LocalStorage().Get(quarantineKey, &kept); err != nil {
		c.t.Fatalf("read back the quarantined snapshot: %v", err)
	}
	if kept.Record.Seed == 0 {
		c.t.Fatal("the failing snapshot was deleted, not quarantined")
	}
	return kept
}
