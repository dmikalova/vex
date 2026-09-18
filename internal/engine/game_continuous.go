package engine

import "fmt"

// This file holds the Continuous effects registry: the flat stack of
// duration-scoped, read-live modifiers a resolving effect installs on the board —
// damage immunity, lost keywords, blanked text, and stat overrides. It is the
// runtime twin of ConstantAbility: where a ConstantAbility is a property of a card
// in play (re-derived from its definition on every read, gone when the card
// leaves play), a ContinuousEffect outlives its source and lasts for a Duration.
// Both feed the same read predicates — Power/armor, hasKeyword, mitigateDamage,
// textBlanked consult the registry and the in-play definitions together.
//
// The registry mirrors the take-control table (game_control.go): a fixed array
// plus count, a generic drop-by-predicate compaction, cleared two ways — by turn
// boundary as each entry's window closes, and never by source (an installed
// effect is not tied to its source staying in play, unlike a ConstantAbility).

// maxContinuous bounds the continuous-effects stack. Each entry is a small flat
// record, so the array stays cheap to snapshot; a 33rd entry panics rather than
// silently dropping an effect.
const maxContinuous = 32

// ContinuousKind names what a continuous effect does to the cards it reaches.
type ContinuousKind uint8

const (
	continuousUnset ContinuousKind = iota
	// ContinuousDamageImmune makes reached creatures unable to be dealt damage.
	ContinuousDamageImmune
	// ContinuousKeywordLost strips the Keywords bitmask from reached creatures.
	ContinuousKeywordLost
	// ContinuousTextBlank blanks reached creatures' text boxes (traits and stats
	// remain).
	ContinuousTextBlank
	// ContinuousStatOverride masks reached creatures' power and/or armor to a fixed
	// value while active — The Pale Star sets each creature to 1 power, 0 armor. It
	// is a read-time mask: the stored counters and armor pool are untouched and
	// revealed again when it lifts.
	ContinuousStatOverride
)

// ContinuousScope says which cards a continuous effect reaches, read live so a
// creature played or gained after the effect resolves is covered too. A compact
// enum rather than a full Target because GameState must stay flat and pointerless
// (ADR 0005) and Target carries a string and an interface.
type ContinuousScope uint8

const (
	scopeUnset ContinuousScope = iota
	// ScopeAllCreatures reaches every creature in play, both sides.
	ScopeAllCreatures
	// ScopeFriendlyCreatures reaches every creature the Controller controls.
	ScopeFriendlyCreatures
	// ScopeEnemyCreatures reaches every creature the Controller's opponent controls.
	ScopeEnemyCreatures
	// ScopeSubject reaches only the single card named by Subject.
	ScopeSubject
)

// ContinuousEffect is one duration-scoped, read-live modifier on the stack. It is
// flat, pointerless, and comparable (ADR 0005): the reached cards are named by a
// Scope enum evaluated live, and the window is two turn numbers rather than a
// closure. ActiveFrom is the first turn it is read as active; the end-of-turn
// sweep drops it once Turn passes ExpiresEnd.
type ContinuousEffect struct {
	Kind       ContinuousKind
	Scope      ContinuousScope
	Subject    LocalID
	Controller int8
	// ActiveFrom is the turn number from which the effect reads as active. An
	// OpponentNextTurn effect is dormant until the opponent's turn, so ActiveFrom is
	// the turn after it is installed.
	ActiveFrom int
	// ExpiresEnd is the turn number at whose end-of-turn seam the effect is dropped.
	ExpiresEnd int
	// Keywords is the lost-keyword bitmask for ContinuousKeywordLost.
	Keywords uint16
	// Power and Armor are the masked values for ContinuousStatOverride.
	Power    int8
	Armor    int8
	HasPower bool
	HasArmor bool
}

// continuousLifetime turns a Duration into the [ActiveFrom, ExpiresEnd] turn-number
// window, read against the current turn. RemainderOfPlayerTurn is live now and
// lifts at the end of this turn; StartOfPlayerNextTurn is live now and lasts
// through the opponent's turn; OpponentNextTurn is dormant now and bites only
// during the opponent's next turn.
func continuousLifetime(d Duration, turn int) (activeFrom, expiresEnd int) {
	switch d {
	case OpponentNextTurn:
		return turn + 1, turn + 1
	case StartOfPlayerNextTurn:
		return turn, turn + 1
	default: // RemainderOfPlayerTurn
		return turn, turn
	}
}

// addContinuous pushes a continuous effect onto the stack, filling in its window
// from the duration and the current turn. An entry past capacity panics — a caught
// invariant, never a silent drop.
func (g *Game) addContinuous(e ContinuousEffect, d Duration) {
	if int(g.State.ContinuousCount) >= maxContinuous {
		panic(fmt.Sprintf("continuous stack full: cannot add %d effect", e.Kind))
	}
	e.ActiveFrom, e.ExpiresEnd = continuousLifetime(d, g.State.Turn)
	g.State.Continuous[g.State.ContinuousCount] = e
	g.State.ContinuousCount++
}

// dropContinuous removes every stack entry the predicate matches, keeping the
// survivors contiguous and order-by-construction (oldest first).
//
// This deliberately does not share a primitive with the lasting registry's
// removeLastingAt. That one deletes a single known index; this one is a bulk
// filter-in-place, and expressing it as repeated index deletes would make it
// quadratic. The two stay distinct operations.
func (g *Game) dropContinuous(match func(ContinuousEffect) bool) {
	w := 0
	for i := 0; i < int(g.State.ContinuousCount); i++ {
		e := g.State.Continuous[i]
		if match(e) {
			continue
		}
		g.State.Continuous[w] = e
		w++
	}
	for i := w; i < int(g.State.ContinuousCount); i++ {
		g.State.Continuous[i] = ContinuousEffect{}
	}
	g.State.ContinuousCount = uint8(w)
}

// clearExpiredContinuous drops every continuous effect whose window has closed. It
// runs at the very end of the turn — after the end-of-turn ability window, just
// before the turn hands over — so an end-of-turn ability still sees this turn's
// effects active, and a "during your opponent's next turn" effect installed this
// turn survives into that turn.
func (g *Game) clearExpiredContinuous() {
	g.dropContinuous(func(e ContinuousEffect) bool { return g.State.Turn >= e.ExpiresEnd })
}

// continuousReaches reports whether an active continuous effect reaches card id,
// resolving its Scope live from the current board.
func (g *Game) continuousReaches(e ContinuousEffect, id LocalID) bool {
	switch e.Scope {
	case ScopeSubject:
		return e.Subject == id
	case ScopeAllCreatures:
		return g.TypeOf(id) == Creature
	case ScopeFriendlyCreatures:
		return g.TypeOf(id) == Creature && g.controller(id) == int(e.Controller)
	default: // ScopeEnemyCreatures
		return g.TypeOf(id) == Creature && g.controller(id) != int(e.Controller)
	}
}

// continuousActive reports whether any continuous effect of the given kind is
// active this turn and reaches card id — the OR-combined read for the monotonic
// kinds (damage immunity, text blank).
func (g *Game) continuousActive(kind ContinuousKind, id LocalID) bool {
	for i := 0; i < int(g.State.ContinuousCount); i++ {
		e := g.State.Continuous[i]
		if e.Kind == kind && g.State.Turn >= e.ActiveFrom && g.continuousReaches(e, id) {
			return true
		}
	}
	return false
}

// DamageImmune reports whether a creature currently cannot be dealt damage through
// a continuous effect (Shield of Justice, Lucky Dice, Protectrix).
func (g *Game) DamageImmune(id LocalID) bool {
	return g.continuousActive(ContinuousDamageImmune, id)
}

// continuousLostKeywords ORs together the lost-keyword bitmasks of every active
// ContinuousKeywordLost effect reaching card id.
func (g *Game) continuousLostKeywords(id LocalID) uint16 {
	var lost uint16
	for i := 0; i < int(g.State.ContinuousCount); i++ {
		e := g.State.Continuous[i]
		if e.Kind == ContinuousKeywordLost && g.State.Turn >= e.ActiveFrom &&
			g.continuousReaches(e, id) {
			lost |= e.Keywords
		}
	}
	return lost
}

// continuousPowerOverride returns the power the latest active stat-override
// masking card id sets, and whether any does. Later entries win (LIFO): the most
// recently installed override is the one in force.
func (g *Game) continuousPowerOverride(id LocalID) (int, bool) {
	power, has := 0, false
	for i := 0; i < int(g.State.ContinuousCount); i++ {
		e := g.State.Continuous[i]
		if e.Kind == ContinuousStatOverride && e.HasPower && g.State.Turn >= e.ActiveFrom &&
			g.continuousReaches(e, id) {
			power, has = int(e.Power), true
		}
	}
	return power, has
}

// continuousArmorOverride returns the armor the latest active stat-override masking
// card id sets, and whether any does. Later entries win (LIFO).
func (g *Game) continuousArmorOverride(id LocalID) (int, bool) {
	armor, has := 0, false
	for i := 0; i < int(g.State.ContinuousCount); i++ {
		e := g.State.Continuous[i]
		if e.Kind == ContinuousStatOverride && e.HasArmor && g.State.Turn >= e.ActiveFrom &&
			g.continuousReaches(e, id) {
			armor, has = int(e.Armor), true
		}
	}
	return armor, has
}
