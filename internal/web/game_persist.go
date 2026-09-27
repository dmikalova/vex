package web

import (
	"math/rand"
	"sort"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/cards"
	"github.com/dmikalova/vex/internal/engine"
	"github.com/dmikalova/vex/internal/match"
	"github.com/dmikalova/vex/internal/session"
)

// This file keeps a match alive across page loads: saving it to local storage,
// restoring it, and — when there is nothing to restore — dealing a new one.

// save writes the current match to local storage. It runs after every action and
// before a hot-reload, so a reload or a later visit resumes from the latest
// state. What is persisted is the session's Record (ADR 0039) — version, seed,
// sets and the ordered command log — from which the state and the typed log are
// replayed. There is no half-resolved action to hold back any more: a prompt's
// answers so far are ordinary recorded commands, and a reload that lands mid-prompt
// rebuilds the prompt along with everything else.
func (g *game) save(ctx app.Context) {
	if g.s == nil {
		return
	}
	g.writeSnapshot(ctx, snapshot{
		Version: snapshotVersion,
		Record:  g.s.Record(),
		UI:      g.savedUI(),
	})
}

// storageFullNotice is the standing message shown once the match can no longer be
// written. It names the cause, because the player's only remedy is outside the
// app: free storage, or leave private browsing.
const storageFullNotice = "This match is no longer being saved — browser storage " +
	"is full or unavailable. Closing this tab will lose it."

// writeSnapshot persists the match, and on a failed write frees the storage the
// client can spare and retries once. The expected failure is a full quota: the
// command log grows with the match, so the write that fails is larger than the
// snapshot already occupying the slot, and dropping that stale snapshot is
// usually enough room for the new one. Only a second failure raises the notice.
func (g *game) writeSnapshot(ctx app.Context, snap snapshot) {
	store := ctx.LocalStorage()
	if err := store.Set(g.matchKey(), snap); err == nil {
		g.clearStorageNotice()
		return
	}
	// The style gallery's scroll memo belongs to another page, and the snapshot in
	// the slot is the one this write replaces, so neither is worth keeping over the
	// match itself.
	store.Del(styleScrollKey)
	store.Del(g.matchKey())
	if err := store.Set(g.matchKey(), snap); err != nil {
		g.setNotice(storageFullNotice)
		return
	}
	g.clearStorageNotice()
}

// clearStorageNotice takes the storage notice down once a write succeeds, leaving
// a notice raised for any other reason alone — the banner holds one message, and a
// successful save says nothing about an unrelated fault.
func (g *game) clearStorageNotice() {
	if g.notice == storageFullNotice {
		g.clearNotice()
	}
}

// savedUI captures the view state to carry across a reload.
func (g *game) savedUI() savedUI {
	ui := savedUI{
		SidebarCollapsed: g.sidebarCollapsed,
		ZonesPlayer:      g.zonesPlayer,
		KeysOpen:         g.keysOpen,
		PickerOpen:       g.pickerOpen,
	}
	// A viewer a prompt opened belongs to that prompt, and the prompt does not
	// survive the reload, so it comes back closed.
	if g.promptZone != "" {
		ui.ZonesPlayer = -1
	}
	return ui
}

// restoreUI puts back the view a reload would otherwise throw away. The viewer's
// player is range-checked because a snapshot is outside data. The picker is only
// ever opened, never closed, because a resume that lands on a name-a-card prompt
// has already opened it to answer that prompt.
func (g *game) restoreUI(ui savedUI) {
	g.sidebarCollapsed = ui.SidebarCollapsed
	g.keysOpen = ui.KeysOpen
	g.pickerOpen = g.pickerOpen || (ui.PickerOpen && g.eng().Manual())
	g.zonesPlayer = -1
	if ui.ZonesPlayer == 0 || ui.ZonesPlayer == 1 {
		g.zonesPlayer = ui.ZonesPlayer
	}
}

// resume rebuilds the match from a saved snapshot, reporting whether it restored
// one. A missing, wrong-version, or engine-incompatible snapshot is dropped and
// resume returns false, so the caller deals a fresh game. The snapshot holds only
// the session's Record, so the resume replays its command log from a fresh deal to
// regenerate the exact state and typed log; a replay that diverges (an older card
// pool, a since-changed action) is caught and started over.
func (g *game) resume(ctx app.Context) (ok bool) {
	store := ctx.LocalStorage()
	if !store.Contains(g.matchKey()) {
		return false
	}
	var snap snapshot
	if err := store.Get(g.matchKey(), &snap); err != nil ||
		snap.Version != snapshotVersion || snap.Record.Seed == 0 {
		store.Del(g.matchKey())
		return false
	}
	// Replaying a log against a different engine/card pool can panic on an
	// out-of-range id; recover and fall back to a fresh deal.
	defer func() {
		if recover() != nil {
			store.Del(g.matchKey())
			ok = false
		}
	}()

	if err := g.replayRecord(snap.Record); err != nil {
		store.Del(g.matchKey())
		return false
	}
	g.restoreUI(snap.UI)
	// A resumed match already holds its whole log, so only lines produced after
	// this point are news; start the toast caught up rather than replaying the
	// entire history into one banner on the first render.
	g.toastSeen = len(g.eng().Log)
	g.status = ""
	return true
}

// replayRecord rebuilds a match from a persisted record by feeding its commands
// back through the same apply path a click takes. session.Load would rebuild the
// state just as faithfully, but it cannot rebuild what is the CLIENT's knowledge
// rather than the log's: where one player action stops and the next begins, which
// is both the undo cursor and the log bubbles. Driving the replay here recovers
// them, because a root boundary is observable — the session is waiting on a
// RequestAction — and a manual force-edit is a root by its kind.
//
// It refuses a record from another session version the way session.Load does:
// there is no forward compatibility, so a stale record is dropped rather than
// misread (ADR 0039).
func (g *game) replayRecord(rec session.Record) error {
	if rec.Version != session.Version {
		return session.ErrVersion
	}
	g.startSession(rec.Seed, rec.Sets)
	for _, cmd := range rec.Commands {
		req, live := g.s.Pending()
		if cmd.Kind.IsManual() || (live && req.Kind == engine.RequestAction) {
			g.beginAction()
		}
		if err := g.applyRecorded(cmd); err != nil {
			return err
		}
	}
	g.settlePending()
	g.afterRebuild()
	return nil
}

// newMatch deals a fresh match. Both sides are driven by the same person
// (hotseat). An injected fixedSeed replaces the clock, so a ui-test scenario deals
// the same cards on every pass.
func (g *game) newMatch() {
	if g.fixedSeed != 0 {
		g.dealMatch(g.fixedSeed)
		return
	}
	g.dealMatch(time.Now().UnixNano())
}

// dealMatch deals a match from a given seed, which fixes the decks and every card
// id in them, and drives it to its first decision: the first-player roll answers
// itself, and the session then stops at the first player's mulligan.
func (g *game) dealMatch(seed int64) {
	g.startSession(seed, g.setNames)
	g.rollFirstPlayer()
	g.settlePending()
	g.afterRebuild()
	// Before the first turn the client is waiting to be asked for a house, whatever
	// the engine is doing in setup, so the resting phase is set after the rebuild
	// rather than derived from an engine that has not started a turn yet.
	g.phase = phaseHouse
}

// startSession stands a fresh session up over the deal seed and sets describe and
// resets everything the client keeps about the match it replaces. The driving
// action is the engine's canonical turn loop, so setup, the mulligans and every
// turn come back as Requests rather than being driven from here.
func (g *game) startSession(seed int64, sets [2]string) {
	g.seed, g.setNames = seed, sets
	g.indexCards()
	g.clearPrompts()
	g.hasAbilityPick = false
	g.resetBadgePreview()
	g.s = session.New(seed, sets, g.dealFor, runMatch, g.lookupCard)
	// No request has been presented for this session yet, and promptAt counts
	// commands, so -1 cannot collide with one.
	g.promptAt = -1
	// The new deal shares no board with the old one, so the log grouping, the undo
	// marks, the redo history and the animation baseline all start over: otherwise
	// the next action diffs against the previous game's cards and flies them all off
	// to the discard, and stale marks bubble the fresh log at the wrong places.
	g.logGroups = nil
	g.rootMarks = nil
	g.redoLog = nil
	g.clearFlashes()
	g.inPlayPrev = g.inPlaySet()
	g.phase = phaseHouse
	g.clearSelection()
	g.zonesPlayer = -1
	g.deckOpen = [2]bool{}
	g.status = ""
}

// runMatch is the action every session the client drives is given: the engine's
// canonical turn loop (ADR 0039, 0040). It owns StartGame, the mulligan prompts
// and the whole sequence of turns, so the client no longer drives a turn of its
// own — it answers the requests the loop suspends on.
func runMatch(eg *engine.Game) { eg.RunMatch() }

// rollFirstPlayer answers the match's one setup decision without asking. Hotseat
// has nobody to ask, so the client flips a coin DERIVED FROM THE SEED: one number
// then reproduces a whole match start, which is what keeps a fixed-seed browser
// scenario stable from pass to pass. It is recorded as an ordinary
// CommandSetFirstPlayer whose decider is engine.RolledFirstPlayer, so the game log
// reads "goes first by random choice" rather than naming a chooser there was none
// of.
func (g *game) rollFirstPlayer() {
	req, live := g.pending()
	if !live || req.Kind != engine.RequestFirstPlayer {
		return
	}
	_ = g.s.Apply(engine.Command{
		Kind:   engine.CommandSetFirstPlayer,
		Player: rand.New(rand.NewSource(g.seed)).Intn(2), //nolint:gosec // a coin flip
		Index:  engine.RolledFirstPlayer,
	})
}

// dealFor is the session's Setup: it deals the match for a seed and set pair and
// returns the engine game. Everything ELSE about a deal that the client draws but
// the engine does not hold — each player's deck houses, which dealt cards are
// Mavericks or Legacy cards, the static rosters the deck-list popover reads — is a
// SIDE EFFECT of this closure, because session.Setup is only shaped to return the
// game. It runs on every replay, which is harmless because the deal is
// deterministic: the same seed and sets always produce the same decks, the same
// card ids and the same emblems.
func (g *game) dealFor(seed int64, sets [2]string) *engine.Game {
	eg, houses, mavericks, legacies, rosters, err := match.NewWithSets(
		"Player 1",
		"Player 2",
		seed,
		sets,
	)
	if err != nil {
		// A set name the deck generator does not know means the picker, or a restored
		// snapshot, is carrying a set this build no longer ships. Say so and deal the
		// default set, so the player gets a playable match instead of an empty board
		// and is told which choice was dropped. The record keeps the name that was
		// entered, so a replay falls back here again and rebuilds the same decks.
		g.setNotice(err.Error() + " — dealing the default set instead.")
		eg, houses, mavericks, legacies, rosters, _ = match.NewWithSets(
			"Player 1",
			"Player 2",
			seed,
			[2]string{},
		)
	}
	// Declare each player's deck houses so a forced house they lack is ignored
	// (cannot overrides must); the human is prompted for each key's colour.
	eg.SetPlayerHouses(0, houses[0])
	eg.SetPlayerHouses(1, houses[1])
	g.deckHouses = houses
	g.rosters = rosters
	g.mavericks = idSet(mavericks)
	g.legacy = idSet(legacies)
	return eg
}

// idSet collapses both players' id lists into one lookup, for the emblems that
// belong to a dealt card rather than to a side.
func idSet(ids [2][]engine.LocalID) map[engine.LocalID]bool {
	set := make(map[engine.LocalID]bool)
	for _, side := range ids {
		for _, id := range side {
			set[id] = true
		}
	}
	return set
}

// lookupCard is the session's ResolveCard: it maps the name a recorded
// CommandManualAddCard carries back to its definition, so a replay re-registers the
// card to the same id the rest of the log refers to.
func (g *game) lookupCard(name string) (engine.CardDefinition, bool) {
	g.indexCards()
	def, ok := g.defByName[name]
	if !ok {
		return engine.CardDefinition{}, false
	}
	return *def, true
}

// indexCards builds the card indexes the client browses and looks names up in: the
// name index behind log mentions and manual adds, and the house-then-name ordered
// pool the manual card picker walks (rather than whichever cards register first).
func (g *game) indexCards() {
	if g.defByName == nil {
		g.defByName = cardsByName()
	}
	if g.allDefs == nil {
		g.allDefs = cards.All()
		sort.Slice(g.allDefs, func(i, j int) bool {
			a, b := g.allDefs[i], g.allDefs[j]
			if a.House != b.House {
				return a.House < b.House
			}
			return a.Name < b.Name
		})
	}
}
