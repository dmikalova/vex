package engine

import "fmt"

// maxCounterEntries bounds the global counter side-table: the number of distinct
// (card, kind) pairs that can carry a generic counter at once. A card rarely
// holds more than one or two kinds and a whole game rarely has more than two
// kinds in play, so 64 is far above any real board (ADR 0024). A 65th distinct
// pair panics rather than silently dropping a counter.
const maxCounterEntries = 64

// CounterKind identifies a generic counter — a card-placed marker that sits on an
// in-play card, stacks, does nothing on its own, and is read only by the cards
// that name it (ADR 0024). Each kind carries a display noun for card text; that
// is its only per-kind data, and it needs no rulebook term of its own.
type CounterKind uint8

const (
	// CounterNone is the unset zero value: no counter.
	CounterNone CounterKind = iota
	// CounterDoom is a doom counter — Wretched Doll destroys every creature
	// carrying one when there is one in play, and otherwise places a fresh one.
	CounterDoom
	// NumCounterKinds is one past the last real kind, so callers can range over
	// CounterNone+1 .. NumCounterKinds to visit every counter (the web icon
	// completeness test does, to force a new kind to ship its own unique icon).
	NumCounterKinds = int(iota)
)

// valid reports whether the kind names a real counter.
func (k CounterKind) valid() bool { return k != CounterNone }

// noun renders the counter's display name for card text, e.g. "doom counter".
func (k CounterKind) noun() string {
	switch k {
	case CounterDoom:
		return "doom counter"
	default:
		return "counter"
	}
}

// CounterEntry is one (card, kind) pair in the global side-table, with its count
// folded into N — a creature carrying two doom counters is one entry with N=2,
// not two entries (ADR 0024). N saturates at 255.
type CounterEntry struct {
	Card LocalID
	Kind CounterKind
	N    uint8
}

// counterIndex returns the position of the (id, kind) entry in the table, or -1
// if the card carries no counter of that kind.
func (g *Game) counterIndex(id LocalID, kind CounterKind) int {
	for i := 0; i < int(g.State.CounterCount); i++ {
		e := g.State.Counters[i]
		if e.Card == id && e.Kind == kind {
			return i
		}
	}
	return -1
}

// PlaceCounter puts n counters of a kind on an in-play card, folding them into
// the card's existing entry (or opening a fresh one), saturating at 255. Placing
// on a card that is not in play, or a non-positive count, is a no-op. A 65th
// distinct (card, kind) pair panics — a caught invariant, never a silent drop.
func (g *Game) PlaceCounter(id LocalID, kind CounterKind, n int) {
	if n <= 0 || g.stateOf(id) == nil {
		return
	}
	if i := g.counterIndex(id, kind); i >= 0 {
		g.State.Counters[i].N = saturateCounter(int(g.State.Counters[i].N) + n)
		return
	}
	if int(g.State.CounterCount) >= maxCounterEntries {
		panic(fmt.Sprintf("counter table full: cannot place %s on card %d", kind.noun(), id))
	}
	g.State.Counters[g.State.CounterCount] = CounterEntry{
		Card: id,
		Kind: kind,
		N:    saturateCounter(n),
	}
	g.State.CounterCount++
}

// CountersOn returns how many counters of a kind sit on a card.
func (g *Game) CountersOn(id LocalID, kind CounterKind) int {
	if i := g.counterIndex(id, kind); i >= 0 {
		return int(g.State.Counters[i].N)
	}
	return 0
}

// clearCounters drops every counter a card carries, compacting the table in
// place. A card sheds its counters when it leaves play; removeFromPlay calls this
// for every exit (ADR 0024).
func (g *Game) clearCounters(id LocalID) {
	for i := 0; i < int(g.State.CounterCount); {
		if g.State.Counters[i].Card == id {
			g.removeCounterEntryAt(i)
			continue
		}
		i++
	}
}

// removeCounterEntryAt deletes the entry at position i, shifting the tail left to
// keep the live entries contiguous and order-by-construction (ADR 0024).
func (g *Game) removeCounterEntryAt(i int) {
	n := int(g.State.CounterCount)
	copy(g.State.Counters[i:], g.State.Counters[i+1:n])
	g.State.CounterCount--
	g.State.Counters[g.State.CounterCount] = CounterEntry{}
}

// saturateCounter clamps a counter total to the uint8 range.
func saturateCounter(n int) uint8 {
	if n > 255 {
		return 255
	}
	return uint8(n)
}
