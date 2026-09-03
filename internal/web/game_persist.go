package web

import (
	"sort"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// This file keeps a match alive across page loads: saving it to local storage,
// restoring it, and — when there is nothing to restore — dealing a new one.

// save writes the current match to local storage. It runs after every action and
// before a hot-reload, so a reload or a later visit resumes from the latest
// state.
func (g *game) save(ctx app.Context) {
	if g.g == nil {
		return
	}
	state, log, groups := g.g.State, g.g.Log, g.logGroups
	// A half-resolved action cannot be persisted: the rest of it lives in a
	// goroutine that a reload kills, so the Æmber it already spent and the card it
	// already took from hand would be lost with no prompt left to finish them.
	// Save the point before it instead, so the reload lands where the player can
	// take the action again.
	if g.busy && len(g.undo) > 0 {
		e := g.undo[len(g.undo)-1]
		state, log, groups = e.state, e.log, e.groups
	}
	saved := make([]savedLine, len(log))
	for i, rec := range log {
		rule, player := ruleOf(rec)
		saved[i] = savedLine{
			Frame:  rec.Frame,
			Text:   rec.Text(g.g),
			Rule:   rule,
			Player: player,
		}
	}
	_ = ctx.LocalStorage().Set(persistKey, snapshot{
		Version: snapshotVersion,
		Seed:    g.seed,
		State:   state,
		Log:     saved,
		Groups:  groups,
		Manual:  g.manualAdds,
	})
}

// resume rebuilds the match from a saved snapshot, reporting whether it restored
// one. A missing, wrong-version, or engine-incompatible snapshot is dropped and
// resume returns false, so the caller deals a fresh game. Reconstruction from the
// seed is deterministic, so the rebuilt catalog and card ids match the saved
// state exactly; any residual mismatch (an older card pool) is caught before use.
func (g *game) resume(ctx app.Context) (ok bool) {
	store := ctx.LocalStorage()
	if !store.Contains(persistKey) {
		return false
	}
	var snap snapshot
	if err := store.Get(persistKey, &snap); err != nil ||
		snap.Version != snapshotVersion || snap.Seed == 0 {
		store.Del(persistKey)
		return false
	}
	// Rebuilding or reading a state from a different engine/card pool can panic on
	// an out-of-range id; recover and fall back to a fresh deal.
	defer func() {
		if recover() != nil {
			store.Del(persistKey)
			ok = false
		}
	}()

	g.seed = snap.Seed
	eg, houses, mavericks := match.NewWithMavericks("Player 1", "Player 2", snap.Seed)
	g.install(eg, houses, mavericks)
	if !g.replayManualAdds(snap.Manual) {
		store.Del(persistKey)
		return false
	}
	g.g.State = snap.State
	if !g.stateReadsCleanly() {
		store.Del(persistKey)
		return false
	}
	// Restore the saved log (a real match always has at least the turn-1 header, so
	// an empty one means a pre-log snapshot — keep what install produced).
	if len(snap.Log) > 0 {
		restored := make([]engine.Record, len(snap.Log))
		for i, line := range snap.Log {
			entry := engine.RestoredEntry{Line: line.Text}
			restored[i] = engine.Record{Frame: line.Frame, Entry: entry}
			if line.Rule != ruleNone {
				restored[i].Entry = restoredRule{
					RestoredEntry: entry,
					Rule:          line.Rule,
					Player:        line.Player,
				}
			}
		}
		g.g.Restore(restored)
		g.logGroups = snap.Groups
	}
	g.settlePhase()
	g.clearSelection()
	g.zonesPlayer = -1
	g.status = ""
	return true
}

// replayManualAdds re-registers the cards manual mode added, in the order they
// were added, so the rebuilt catalog hands out the same ids the saved state
// refers to. It reports false when a name no longer exists in the card pool,
// which leaves the ids misaligned and the snapshot unusable.
func (g *game) replayManualAdds(adds []manualAdd) bool {
	for _, a := range adds {
		def, ok := g.defByName[a.Name]
		if !ok {
			return false
		}
		g.g.Register(*def, a.Player)
	}
	g.manualAdds = adds
	return true
}

// stateReadsCleanly reports whether every card id in the restored state resolves
// against the rebuilt catalog. A snapshot from a different card pool can hold ids
// out of range, which panics on lookup; recovering turns that into a clean
// start-over signal.
func (g *game) stateReadsCleanly() (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	for p := 0; p < 2; p++ {
		for _, zone := range [][]engine.LocalID{
			g.g.Hand(p), g.g.Deck(p), g.g.Battleline(p), g.g.Artifacts(p),
			g.g.Discard(p), g.g.Archives(p), g.g.Purge(p),
		} {
			for _, id := range zone {
				g.g.Def(id)
			}
		}
	}
	return true
}

// newMatch seeds a new game, wires the shared human chooser to both players, and
// deals random decks. Both sides are driven by the same person (hotseat).
func (g *game) newMatch() {
	g.seed = time.Now().UnixNano()
	eg, houses, mavericks := match.NewWithMavericks("Player 1", "Player 2", g.seed)
	g.install(eg, houses, mavericks)
	// Clear the previous game's log grouping and undo/redo history. newMatch resets
	// the engine log to a single turn-1 header, so stale marks (with larger Start
	// indices from the old, longer log) would bubble the fresh log at the wrong
	// places.
	g.logGroups = nil
	g.undo = nil
	g.redo = nil
	g.manualAdds = nil
	g.g.StartTurn(0) // no Æmber yet, so the opening forge step is a no-op
	// The new deal shares no board with the old one, so the animation baseline is
	// reset too: otherwise the next action diffs against the previous game's cards
	// and flies them all off to the discard.
	g.clearFlashes()
	g.inPlayPrev = g.inPlaySet()
	g.phase = phaseHouse
	g.clearSelection()
	g.zonesPlayer = -1
	g.status = ""
}

// install wires a freshly built engine game into the component: it attaches the
// shared human chooser to both players and records the harness, its deck houses,
// and which dealt cards are Mavericks. The caller sets the starting phase.
func (g *game) install(
	eg *engine.Game,
	houses [2][]engine.House,
	mavericks [2][]engine.LocalID,
) {
	ch := &webChooser{g: g, reply: make(chan chooseReply, 1), optionReply: make(chan int, 1)}
	eg.SetChooser(0, ch)
	eg.SetChooser(1, ch)
	// Declare each player's deck houses so a forced house they lack is ignored
	// (cannot overrides must); the human is prompted for each key's colour.
	eg.SetPlayerHouses(0, houses[0])
	eg.SetPlayerHouses(1, houses[1])
	g.g = eg
	g.chooser = ch
	g.deckHouses = houses
	g.mavericks = make(map[engine.LocalID]bool)
	for _, ids := range mavericks {
		for _, id := range ids {
			g.mavericks[id] = true
		}
	}
	if g.defByName == nil {
		g.defByName = cardsByName()
	}
	if g.allDefs == nil {
		g.allDefs = cards.All()
		// Sort by house then name so the manual card picker browses every house in
		// order rather than showing whichever cards register first.
		sort.Slice(g.allDefs, func(i, j int) bool {
			a, b := g.allDefs[i], g.allDefs[j]
			if a.House != b.House {
				return a.House < b.House
			}
			return a.Name < b.Name
		})
	}
}
