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
// state. Only the command log is persisted (ADR 0039); the state and typed log a
// resume needs are replayed from it.
func (g *game) save(ctx app.Context) {
	if g.g == nil {
		return
	}
	inputs := g.inputs
	// A half-resolved action cannot be persisted: the rest of it lives in a
	// goroutine that a reload kills, so a choice it has not recorded yet would be
	// lost with no prompt left to make it. Persist only the committed prefix —
	// every input up to the last completed root — so the reload lands where the
	// player can take the pending action again.
	if g.busy && len(g.rootMarks) > 0 {
		inputs = inputs[:g.rootMarks[len(g.rootMarks)-1]]
	} else if g.busy {
		inputs = nil
	}
	_ = ctx.LocalStorage().Set(persistKey, snapshot{
		Version:  snapshotVersion,
		Seed:     g.seed,
		SetNames: g.setNames,
		Inputs:   append([]input(nil), inputs...),
		UI:       g.savedUI(),
	})
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
// player is range-checked because a snapshot is outside data.
func (g *game) restoreUI(ui savedUI) {
	g.sidebarCollapsed = ui.SidebarCollapsed
	g.keysOpen = ui.KeysOpen
	g.pickerOpen = ui.PickerOpen && g.g.Manual()
	g.zonesPlayer = -1
	if ui.ZonesPlayer == 0 || ui.ZonesPlayer == 1 {
		g.zonesPlayer = ui.ZonesPlayer
	}
}

// resume rebuilds the match from a saved snapshot, reporting whether it restored
// one. A missing, wrong-version, or engine-incompatible snapshot is dropped and
// resume returns false, so the caller deals a fresh game. The snapshot holds only
// the seed, sets, and command log, so the resume replays the log from a fresh
// deal to regenerate the exact state and typed log; a replay that diverges (an
// older card pool, a since-changed action) is caught and started over.
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
	// Replaying a log against a different engine/card pool can panic on an
	// out-of-range id; recover and fall back to a fresh deal.
	defer func() {
		if recover() != nil {
			store.Del(persistKey)
			ok = false
		}
	}()

	g.seed = snap.Seed
	g.setNames = snap.SetNames
	g.inputs = snap.Inputs
	g.recomputeRootMarks()
	if !g.rebuildFromLog() {
		store.Del(persistKey)
		return false
	}
	g.redoLog = nil
	g.settlePhase()
	g.clearSelection()
	g.inPlayPrev = g.inPlaySet()
	g.prevValid = false
	g.restoreUI(snap.UI)
	// A resumed match already holds its whole log, so only lines produced after
	// this point are news; start the toast caught up rather than replaying the
	// entire history into one banner on the first render.
	g.toastSeen = len(g.g.Log)
	g.status = ""
	return true
}

// newMatch seeds a new game, wires the shared human chooser to both players, and
// deals random decks. Both sides are driven by the same person (hotseat).
func (g *game) newMatch() { g.dealMatch(time.Now().UnixNano()) }

// dealMatch deals a match from a given seed, which fixes the decks and every card
// id in them.
func (g *game) dealMatch(seed int64) {
	g.seed = seed
	// A new deal may abandon a prompt the previous match left blocked on its chooser
	// (a mulligan, or a mid-action pick): release that goroutine so it returns
	// rather than leaking, and so its stale UI dispatches no-op against this deal.
	if g.chooser != nil {
		g.chooser.drain()
	}
	// The drained goroutine's own prompt teardown is stale-guarded, so it cannot
	// clear the prompt it left up; clear it here so this deal starts with no stale
	// prompt (a leftover choosingOption would make the next await match it).
	g.clearPrompts()
	eg, houses, mavericks, legacies, rosters := match.NewWithSets(
		"Player 1",
		"Player 2",
		g.seed,
		g.setNames,
	)
	g.install(eg, houses, mavericks, legacies, rosters)
	// Clear the previous game's log grouping and redo history. newMatch resets
	// the engine log to a single turn-1 header, so stale marks (with larger Start
	// indices from the old, longer log) would bubble the fresh log at the wrong
	// places.
	g.logGroups = nil
	g.redoLog = nil
	// A fresh deal starts a fresh command log: the mulligan answers StartGame is
	// about to prompt for are its first recorded inputs.
	g.inputs = nil
	g.rootMarks = nil
	// The new deal shares no board with the old one, so the animation baseline is
	// reset too: otherwise the next action diffs against the previous game's cards
	// and flies them all off to the discard.
	g.clearFlashes()
	g.inPlayPrev = g.inPlaySet()
	g.phase = phaseHouse
	g.clearSelection()
	g.zonesPlayer = -1
	g.deckOpen = [2]bool{}
	g.status = ""
	// Run setup on a background goroutine so the interactive chooser can offer each
	// player their one mulligan: StartGame deals both opening hands, prompts each
	// player in turn (blocking this goroutine on the chooser, never the UI one),
	// then begins the first turn and pauses at the house choice. The finish hops
	// back to the UI goroutine through g.dispatch.
	g.busy = true
	// Capture this deal's engine and chooser so an abandoned setup (a new game
	// started from the mulligan prompt) does not drive the game that replaced it:
	// the goroutine runs its own engine, and its finish no-ops once a newer deal
	// has installed a different chooser.
	ch := g.chooser
	go func() {
		eg.StartGame(0)
		g.dispatch(func(ctx app.Context) {
			if g.chooser != ch {
				return
			}
			g.busy = false
			g.clearFlashes()
			g.inPlayPrev = g.inPlaySet()
			g.phase = phaseHouse
			g.clearSelection()
			g.save(ctx)
		})
	}()
}

// install wires a freshly built engine game into the component: it attaches the
// shared human chooser to both players and records the harness, its deck houses,
// and which dealt cards are Mavericks or Legacy cards. The caller sets the
// starting phase.
func (g *game) install(
	eg *engine.Game,
	houses [2][]engine.House,
	mavericks [2][]engine.LocalID,
	legacies [2][]engine.LocalID,
	rosters [2]match.Roster,
) {
	ch := &webChooser{
		g:             g,
		reply:         make(chan chooseReply, 1),
		optionReply:   make(chan int, 1),
		positionReply: make(chan int, 1),
		cancel:        make(chan struct{}),
	}
	eg.SetChooser(0, ch)
	eg.SetChooser(1, ch)
	// Declare each player's deck houses so a forced house they lack is ignored
	// (cannot overrides must); the human is prompted for each key's colour.
	eg.SetPlayerHouses(0, houses[0])
	eg.SetPlayerHouses(1, houses[1])
	g.g = eg
	g.chooser = ch
	g.deckHouses = houses
	g.rosters = rosters
	g.mavericks = make(map[engine.LocalID]bool)
	for _, ids := range mavericks {
		for _, id := range ids {
			g.mavericks[id] = true
		}
	}
	g.legacy = make(map[engine.LocalID]bool)
	for _, ids := range legacies {
		for _, id := range ids {
			g.legacy[id] = true
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
