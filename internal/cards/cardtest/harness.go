// Package cardtest is a game-like test harness for card-set packages. A test
// declares a whole scenario up front with ct.Play — each player's board, hand,
// house, and pools — then drives it the way the game is played: playing and
// using cards, answering prompts by clicking, and asserting the result with
// h.Expect. Every set's _test.go files import it so this setup is written once
// rather than copied per set.
package cardtest

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Setup declares a whole scenario up front: each player's board, hand, and pools.
// It is the game-like counterpart to wiring a game together call by call — the
// harness reads it once and hands back a ready-to-play match.
type Setup struct {
	P1   Side
	P2   Side
	Seed int64 // 0 uses a fixed default seed for reproducible games
}

// Side is one player's starting position: the house they play this turn, the
// cards in each zone, and their Æmber and keys. Every field is optional; an empty
// Side is a player with an empty board.
type Side struct {
	House    engine.House
	InPlay   []Entry // creatures and artifacts, routed by card type
	Hand     []Entry
	Deck     []Entry
	Discard  []Entry
	Archives []Entry
	Amber    int
	Keys     int
}

// Harness is a running scenario. It owns the game and the two players and bridges
// the engine's choices to interactive clicks. Obtain one from Play.
type Harness struct {
	P1 *Player
	P2 *Player

	// t is testing.TB rather than *testing.T so the harness's own failure paths
	// can be driven by a recording double; a card test always passes its *testing.T.
	t testing.TB
	g *engine.Game

	prompt  chan promptReq
	done    chan actionResult
	current *promptReq
	running bool
	label   string

	// orderScript[p] is the resolution order p wants for their next multi-target
	// ordering, set by Player.Order and consumed once by the bridge Orderer.
	orderScript [2][]any
}

// Player drives one side of a scenario: playing and using cards, ending the turn,
// and answering prompts. P1 is engine player 0 and P2 is engine player 1.
type Player struct {
	h     *Harness
	index int
}

// Play builds a game from a Setup and returns the running harness. Player 1 is
// the active player with their house chosen, so their cards can be played
// immediately. A t.Cleanup guards that the test leaves no prompt unanswered.
func Play(t testing.TB, s Setup) *Harness {
	t.Helper()
	seed := s.Seed
	if seed == 0 {
		seed = 1
	}
	g := engine.NewGame("P1", "P2", seed)
	h := &Harness{
		t:      t,
		g:      g,
		prompt: make(chan promptReq),
		done:   make(chan actionResult),
	}
	h.P1 = &Player{h: h, index: 0}
	h.P2 = &Player{h: h, index: 1}

	// Choose the active player's house before placing cards, so the archives
	// offer (which fires on ChooseHouse) sees an empty archive and never prompts
	// while the default chooser is still installed.
	house := s.P1.House
	if house == engine.HouseNone {
		house = DefaultHouse
	}
	// A made-up house is accepted by the engine but matches no card, so the
	// scenario would silently be one where nothing is playable. Catch it here
	// rather than let the test read a board that never moved.
	if int(house) >= engine.NumHouses {
		t.Fatalf("cardtest: Setup names house %d, which is not one of the %d houses",
			house, engine.NumHouses)
	}
	g.StartTurn(0)
	// ChooseHouse refuses only a locked or forced house, and the first turn of a
	// scenario has neither, so a real house is always accepted here.
	_ = g.ChooseHouse(0, house)

	h.placeSide(0, s.P1)
	h.placeSide(1, s.P2)

	// Install the interactive bridge only now that setup is done, so no engine
	// choice during placement can block on a click.
	g.SetChooser(0, bridgeChooser{h: h, player: 0})
	g.SetChooser(1, bridgeChooser{h: h, player: 1})

	t.Cleanup(h.checkReady)
	return h
}

// Game returns the underlying engine game for assertions the fluent API does not
// yet cover.
func (h *Harness) Game() *engine.Game { return h.g }

// Expect begins a fluent set of assertions about a card, referenced by handle or
// by its definition (resolved to its single instance in the match).
func (h *Harness) Expect(card any) CardExpect {
	h.t.Helper()
	id := h.resolve(card, h.allIDs(), "Expect")
	return CardExpect{h: h, c: Card{h: h, id: id, set: true}}
}

// --- player verbs ---

// Play plays a card from this player's hand, dispatching by its type. The card is
// referenced by handle or definition; its hand position is found automatically.
func (p *Player) Play(card any) {
	p.h.t.Helper()
	idx, id := p.h.findInHand(p.index, card)
	def := p.h.g.Def(id)
	name := def.Name
	switch def.Type {
	case engine.Creature:
		p.h.run(
			"Play "+name,
			func() error { _, err := p.h.g.PlayCreature(p.index, idx, false); return err },
		)
	case engine.Artifact:
		p.h.run(
			"Play "+name,
			func() error { _, err := p.h.g.PlayArtifact(p.index, idx); return err },
		)
	case engine.Tactic:
		p.h.run("Play "+name, func() error { return p.h.g.PlayAction(p.index, idx) })
	case engine.Upgrade:
		p.h.run(
			"Play "+name,
			func() error { _, err := p.h.g.PlayUpgrade(p.index, idx); return err },
		)
	default:
		p.h.t.Fatalf("Play: %s has an unknown type %q", name, def.Type)
	}
}

// Discard discards a card from this player's hand, the turn action that sits
// alongside playing and using a card.
func (p *Player) Discard(card any) {
	p.h.t.Helper()
	_, id := p.h.findInHand(p.index, card)
	p.h.run("Discard "+p.h.g.Name(id), func() error {
		p.h.g.DiscardCardFromHand(p.index, id)
		return nil
	})
}

// ExpectCannotPlay asserts that a card in this player's hand may not be played
// right now — the assertion for cards that gate their own play, such as Kelifi
// Dragon's Æmber threshold.
func (p *Player) ExpectCannotPlay(card any) {
	p.h.t.Helper()
	_, id := p.h.findInHand(p.index, card)
	if err := p.h.g.CanPlay(p.index, id); err == nil {
		p.h.t.Fatalf("%s may be played, want it blocked", p.h.g.Name(id))
	}
}

// ExpectCannotUse asserts that one of this player's cards may not be used right
// now — the assertion for cards that gate their own use, such as Giant Sloth.
func (p *Player) ExpectCannotUse(card any) {
	p.h.t.Helper()
	id := p.h.resolve(card, p.h.inPlayIDs(), "ExpectCannotUse")
	if err := p.h.g.CanUse(p.index, id); err == nil {
		p.h.t.Fatalf("%s may be used, want it blocked", p.h.g.Name(id))
	}
}

// ExpectCannotUseTo asserts that one of this player's cards may not be used one
// specific way while the others stay open — Tireless Crocag cannot reap.
func (p *Player) ExpectCannotUseTo(card any, kind engine.UseKind) {
	p.h.t.Helper()
	id := p.h.resolve(card, p.h.inPlayIDs(), "ExpectCannotUseTo")
	if err := p.h.g.CanUseTo(p.index, id, kind); err == nil {
		p.h.t.Fatalf("%s may be used to %v, want it blocked", p.h.g.Name(id), kind)
	}
}

// Reap reaps with one of this player's creatures.
func (p *Player) Reap(card any) {
	p.h.t.Helper()
	id := p.h.resolve(card, p.h.inPlayIDs(), "Reap")
	p.h.run("Reap "+p.h.g.Name(id), func() error { return p.h.g.Reap(p.index, id) })
}

// Fight has one of this player's creatures fight an enemy creature.
func (p *Player) Fight(attacker, defender any) {
	p.h.t.Helper()
	a := p.h.resolve(attacker, p.h.inPlayIDs(), "Fight")
	d := p.h.resolve(defender, p.h.inPlayIDs(), "Fight")
	p.h.run("Fight", func() error { return p.h.g.Fight(p.index, a, d) })
}

// UseAction fires the "Action:" ability of one of this player's cards.
func (p *Player) UseAction(card any) {
	p.h.t.Helper()
	id := p.h.resolve(card, p.h.inPlayIDs(), "UseAction")
	p.h.run("UseAction "+p.h.g.Name(id), func() error { return p.h.g.UseAction(p.index, id) })
}

// EndTurn ends this player's turn and starts the other player's, leaving them to
// choose a house.
func (p *Player) EndTurn() {
	p.h.t.Helper()
	other := 1 - p.index
	p.h.run("EndTurn", func() error {
		p.h.g.EndPlayPhase(p.index)
		p.h.g.StartTurn(other)
		return nil
	})
}

// ChooseHouse sets this player's active house for the turn.
func (p *Player) ChooseHouse(house engine.House) {
	p.h.t.Helper()
	p.h.run("ChooseHouse", func() error { return p.h.g.ChooseHouse(p.index, house) })
}

// ClickCard answers the pending creature prompt for this player.
func (p *Player) ClickCard(target any) { p.h.t.Helper(); p.h.click(p.index, target) }

// ClickOption answers the pending option prompt for this player by label.
func (p *Player) ClickOption(label string) { p.h.t.Helper(); p.h.clickOption(p.index, label) }

// ClickDone declines the pending optional card prompt for this player — the Done
// button a "you may" or "up to N" choice offers alongside the clickable cards.
func (p *Player) ClickDone() { p.h.t.Helper(); p.h.clickDone(p.index) }

// Order sets the resolution order for this player's next multi-target ordering
// (e.g. which creature an area effect damages first). Cards are matched by handle
// or definition; any not listed keep the engine's order after the listed ones. It
// applies to the next ordering only.
func (p *Player) Order(cards ...any) { p.h.orderScript[p.index] = cards }

// ExpectPrompt asserts the pending prompt's text for this player and returns it
// so the raising card can also be checked with Source.
func (p *Player) ExpectPrompt(text string) Prompt {
	p.h.t.Helper()
	req := p.h.requirePrompt(p.index)
	if req.text != text {
		p.h.t.Errorf("prompt for %s = %q, want %q", playerName(p.index), req.text, text)
	}
	return Prompt{h: p.h, req: req}
}

// Prompt is a captured pending prompt, returned by ExpectPrompt so its source can
// be asserted too.
type Prompt struct {
	h   *Harness
	req promptReq
}

// Source asserts which card raised the prompt, e.g.
// ExpectPrompt("Choose a friendly creature").Source("Anger").
func (p Prompt) Source(name string) Prompt {
	p.h.t.Helper()
	if p.req.source != name {
		p.h.t.Errorf("prompt %q source = %q, want %q", p.req.text, p.req.source, name)
	}
	return p
}

// Amber returns this player's Æmber pool.
func (p *Player) Amber() int { return p.h.g.Aember(p.index) }

// Keys returns this player's forged key count.
func (p *Player) Keys() int { return p.h.g.Keys(p.index) }

// ExpectAmber asserts this player's Æmber pool.
func (p *Player) ExpectAmber(n int) {
	p.h.t.Helper()
	if got := p.Amber(); got != n {
		p.h.t.Errorf("%s Æmber = %d, want %d", playerName(p.index), got, n)
	}
}

// ExpectKeys asserts this player's key count.
func (p *Player) ExpectKeys(n int) {
	p.h.t.Helper()
	if got := p.Keys(); got != n {
		p.h.t.Errorf("%s keys = %d, want %d", playerName(p.index), got, n)
	}
}

// --- placement ---

// placeSide fills a player's zones and pools from a Side.
func (h *Harness) placeSide(player int, s Side) {
	for _, e := range s.InPlay {
		h.placeInPlay(player, e)
	}
	for _, e := range s.Hand {
		h.bind(e, h.g.AddToHand(e.def, player))
	}
	for _, e := range s.Deck {
		h.bind(e, h.g.AddToDeck(e.def, player))
	}
	for _, e := range s.Discard {
		h.bind(e, h.g.AddToDiscard(e.def, player))
	}
	for _, e := range s.Archives {
		h.bind(e, h.g.AddToArchives(e.def, player))
	}
	if s.Amber != 0 {
		h.g.State.Aember[player] = s.Amber
	}
	if s.Keys != 0 {
		h.g.State.Keys[player] = s.Keys
	}
}

// placeInPlay puts a creature on the battleline or an artifact in the artifact
// row, then binds the handle and attaches any upgrades.
func (h *Harness) placeInPlay(player int, e Entry) {
	var id engine.LocalID
	if e.def.Type == engine.Artifact {
		id = h.g.AddArtifact(e.def, player)
	} else {
		id = h.g.AddToBattleline(e.def, player)
	}
	h.bind(e, id)
	for _, u := range e.upgrades {
		up := h.attach(id, u.def)
		if u.bind != nil {
			*u.bind = up
		}
	}
}

// bind fills an entry's handle (if any) with the placed card.
func (h *Harness) bind(e Entry, id engine.LocalID) {
	if e.bind != nil {
		*e.bind = Card{h: h, id: id, set: true}
	}
}

// attach attaches an upgrade to a host creature and refreshes its armor.
func (h *Harness) attach(host engine.LocalID, up engine.CardDefinition) Card {
	id := h.g.Register(up, h.ownerOf(host))
	h.g.AttachUpgrade(host, id)
	h.g.State.Cards[host].ArmorRemaining = int16(h.g.Armor(host))
	return Card{h: h, id: id, set: true}
}

// --- lookup ---

// findInHand resolves a card in a player's hand to its index and id.
func (h *Harness) findInHand(player int, card any) (int, engine.LocalID) {
	h.t.Helper()
	hand := h.g.Hand(player)
	switch v := card.(type) {
	case Card:
		v.require()
		for i, id := range hand {
			if id == v.id {
				return i, id
			}
		}
		h.t.Fatalf("Play: %s is not in %s's hand", h.g.Name(v.id), playerName(player))
		return 0, 0
	case engine.CardDefinition:
		idx, id, n := -1, engine.LocalID(0), 0
		for i, cid := range hand {
			if h.g.Name(cid) == v.Name {
				idx, id, n = i, cid, n+1
			}
		}
		if n == 1 {
			return idx, id
		}
		if n > 1 {
			h.t.Fatalf(
				"Play(%s): %d copies in %s's hand — use ct.Bind to name one",
				v.Name,
				n,
				playerName(player),
			)
		}
		h.t.Fatalf("Play(%s): not in %s's hand", v.Name, playerName(player))
		return 0, 0
	default:
		h.t.Fatalf("Play: cannot play a %T", card)
		return 0, 0
	}
}

// resolve resolves a def-or-handle to a single id within the given id set. A
// handle resolves directly; a definition must match exactly one id.
func (h *Harness) resolve(card any, ids []engine.LocalID, ctx string) engine.LocalID {
	h.t.Helper()
	switch v := card.(type) {
	case Card:
		v.require()
		return v.id
	case engine.CardDefinition:
		found, n := engine.LocalID(0), 0
		for _, id := range ids {
			if h.g.Name(id) == v.Name {
				found, n = id, n+1
			}
		}
		if n == 1 {
			return found
		}
		if n > 1 {
			h.t.Fatalf("%s(%s): %d matching cards — use ct.Bind to name one", ctx, v.Name, n)
		}
		h.t.Fatalf("%s(%s): no matching card found", ctx, v.Name)
		return 0
	default:
		h.t.Fatalf("%s: cannot resolve a %T", ctx, card)
		return 0
	}
}

// ownerOf reports which player currently holds a card, scanning play and hand.
func (h *Harness) ownerOf(id engine.LocalID) int {
	for p := 0; p < 2; p++ {
		if containsID(h.g.Battleline(p), id) || containsID(h.g.Artifacts(p), id) ||
			containsID(h.g.Hand(p), id) {
			return p
		}
	}
	return 0
}

// location reports which zone a card is in.
func (h *Harness) location(id engine.LocalID) Zone {
	for p := 0; p < 2; p++ {
		if containsID(h.g.Battleline(p), id) || containsID(h.g.Artifacts(p), id) {
			return PlayArea
		}
	}
	for p := 0; p < 2; p++ {
		if containsID(h.g.Hand(p), id) {
			return Hand
		}
	}
	for p := 0; p < 2; p++ {
		if containsID(h.g.Discard(p), id) {
			return Discard
		}
	}
	for p := 0; p < 2; p++ {
		if containsID(h.g.Archives(p), id) {
			return Archives
		}
	}
	for p := 0; p < 2; p++ {
		if containsID(h.g.Purge(p), id) {
			return Purge
		}
	}
	for p := 0; p < 2; p++ {
		d := &h.g.State.Deck[p]
		for i := 0; i < int(d.Count); i++ {
			if d.IDs[i] == id {
				return Deck
			}
		}
	}
	for p := 0; p < 2; p++ {
		for _, cid := range h.g.Battleline(p) {
			if containsID(h.g.Upgrades(cid), id) {
				return Attached
			}
		}
	}
	return Gone
}

// inPlayIDs returns every creature and artifact in play, both players.
func (h *Harness) inPlayIDs() []engine.LocalID {
	var out []engine.LocalID
	for p := 0; p < 2; p++ {
		out = append(out, h.g.Battleline(p)...)
		out = append(out, h.g.Artifacts(p)...)
	}
	return out
}

// allIDs returns every card in every zone, plus attached upgrades, both players.
func (h *Harness) allIDs() []engine.LocalID {
	var out []engine.LocalID
	for p := 0; p < 2; p++ {
		out = append(out, h.g.Battleline(p)...)
		out = append(out, h.g.Artifacts(p)...)
		out = append(out, h.g.Hand(p)...)
		out = append(out, h.g.Discard(p)...)
		out = append(out, h.g.Archives(p)...)
		d := &h.g.State.Deck[p]
		for i := 0; i < int(d.Count); i++ {
			out = append(out, d.IDs[i])
		}
		for _, cid := range h.g.Battleline(p) {
			out = append(out, h.g.Upgrades(cid)...)
		}
	}
	return out
}

// nameList renders ids as a bracketed, comma-separated list of names.
func (h *Harness) nameList(ids []engine.LocalID) string {
	names := make([]string, len(ids))
	for i, id := range ids {
		names[i] = h.g.Name(id)
	}
	return "[" + strings.Join(names, ", ") + "]"
}

// describe renders a target (handle or definition) for a message.
func (h *Harness) describe(target any) string {
	switch v := target.(type) {
	case Card:
		if v.set {
			return h.g.Name(v.id)
		}
		return "unbound card"
	case engine.CardDefinition:
		return v.Name
	default:
		return "?"
	}
}

// containsID reports whether ids contains id.
func containsID(ids []engine.LocalID, id engine.LocalID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// reorder arranges ids to match a script of def-or-handle cards; ids not named in
// the script keep the engine's order after the named ones.
func (h *Harness) reorder(ids []engine.LocalID, script []any) []engine.LocalID {
	remaining := append([]engine.LocalID(nil), ids...)
	out := make([]engine.LocalID, 0, len(ids))
	for _, want := range script {
		for i, id := range remaining {
			if h.matchesCard(want, id) {
				out = append(out, id)
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	return append(out, remaining...)
}

// matchesCard reports whether a def-or-handle names the given id.
func (h *Harness) matchesCard(want any, id engine.LocalID) bool {
	switch v := want.(type) {
	case Card:
		return v.set && v.id == id
	case engine.CardDefinition:
		return h.g.Name(id) == v.Name
	default:
		return false
	}
}

// playerName maps a player index to its scenario name.
func playerName(index int) string {
	if index == 0 {
		return "P1"
	}
	return "P2"
}
