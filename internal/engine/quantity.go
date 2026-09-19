package engine

import "fmt"

// Quantity is how many cards a movement verb takes from a zone, and how that
// number reads in the printed text: Takes N, UpTo N (the controller may stop
// early), or AnyNumber (until they decline). It carries both the loop bound and
// the text fragment, so "up to 2 cards" can never desync from a loop that takes
// three (ADR 0006).
//
// It replaces the Amount-plus-AnyNumber-plus-Count trio the movement verbs each
// grew separately, which could spell combinations no card means — an AnyNumber
// beside a fixed count, an Amount beside an Each Selection. Those were rejected in
// validate; here they cannot be written. Where Selection says WHICH cards a verb
// takes, Quantity says HOW MANY, and the two vary independently.
//
// Named Quantity rather than Count, because Count is already this package's "a
// number read off the board" and a Quantity is built FROM one: Takes{Fixed(2)}
// puts its number in the object phrase ("2 cards"), while Takes of a board count
// leads the sentence with a "for each ..." clause instead and keeps the object
// singular.
type Quantity interface {
	// picks is how many times the Selection is asked, and whether that is a bound
	// at all. An unbounded quantity stops only when a pick comes back empty; a
	// bounded one of zero asks nothing, which is not the same thing — a "for each
	// friendly Shard" with no Shard in play shuffles no cards
	// (TestShuffleFromDiscardCount).
	picks(ctx *EffectContext) (n int, bounded bool)
	// object renders the noun phrase for this many cards. noun is the Selection's
	// bare noun ("Sanctum creature") and single its own one-card phrase ("a Sanctum
	// creature"), which a quantity of one returns unchanged.
	object(noun, single string) string
	// leadIn is the board count the sentence opens with ("for each friendly Shard,
	// …"), or nil when the number already reads in the object phrase.
	leadIn() Count
	// fixed is the number when it is known before the game starts — a constant. A
	// board count and an AnyNumber are not known then, so they report false. Text
	// renders without a game state, so every text-time decision reads this rather
	// than picks.
	fixed() (n int, ok bool)
	// optional reports that the controller may stop short of the number. It is what
	// makes a prompt declinable for a verb with no Selection to carry the choice
	// (PutChosen).
	optional() bool
	// validate rejects a quantity whose own number was left unset or is not a
	// number of cards, so a half-written Takes{} cannot pass for "one card".
	validate() error
}

// quantityValidate is q.validate, nil-safe. An omitted Quantity is the game's own
// singular default ("archive a card"); a written one must say a real number.
func quantityValidate(q Quantity) error {
	if q == nil {
		return nil
	}
	return q.validate()
}

// validateCardCount rejects a count that cannot be a number of cards: unset, or a
// constant below one.
func validateCardCount(name string, n Count) error {
	if n == nil {
		return fmt.Errorf("%s: N must be set", name)
	}
	if fixed, ok := countFixed(n); ok {
		return positiveCount(name, "N", fixed)
	}
	return nil
}

// quantityPicks is q.picks with the nil default of a single pick, so a verb that
// leaves Quantity unset still takes one card.
func quantityPicks(q Quantity, ctx *EffectContext) (n int, bounded bool) {
	if q == nil {
		return 1, true
	}
	return q.picks(ctx)
}

// quantityObject is q.object with the nil default of the Selection's own one-card
// phrase.
func quantityObject(q Quantity, noun, single string) string {
	if q == nil {
		return single
	}
	return q.object(noun, single)
}

// quantityLeadIn is q.leadIn, nil-safe.
func quantityLeadIn(q Quantity) Count {
	if q == nil {
		return nil
	}
	return q.leadIn()
}

// quantitySingle reports that the quantity is known before the game to be exactly
// one card, which is what lets a verb render as a single clickable "you may" and
// answer it with one click (PurgeCard.declinable).
func quantitySingle(q Quantity) bool {
	n, ok := quantityFixed(q)
	return ok && n == 1
}

// quantityFixed is q.fixed, nil-safe. An unset quantity is the constant one.
func quantityFixed(q Quantity) (int, bool) {
	if q == nil {
		return 1, true
	}
	return q.fixed()
}

// FixedCardCount is how many cards a Quantity takes when that is a constant above
// one, and 0 otherwise. It is the numeral the web renderer badges a zone glyph
// with, so a board-scaled or single-card quantity shows a bare glyph.
func FixedCardCount(q Quantity) int {
	if n, ok := quantityFixed(q); ok && n > 1 {
		return n
	}
	return 0
}

// quantityOptional is q.optional, nil-safe. An unset quantity takes its one card
// without asking to stop.
func quantityOptional(q Quantity) bool { return q != nil && q.optional() }

// countFixed is a Count's value when it is a constant, which is what text may read
// before there is a game state.
func countFixed(n Count) (int, bool) {
	f, ok := n.(Fixed)
	return int(f), ok
}

// countedObject renders a Count's worth of noun. A constant count prints its
// number ("2 cards"); a board count keeps the singular phrase, because its number
// reads as the sentence's leading "for each …" clause instead.
func countedObject(n Count, noun, single string) string {
	fixed, ok := countFixed(n)
	if !ok || fixed <= 1 {
		return single
	}
	return countNoun(fixed, noun)
}

// countedLeadIn is the Count a sentence must open with: a board count, whose
// number has nowhere else to go. A constant prints in the object instead.
func countedLeadIn(n Count) Count {
	if _, fixed := countFixed(n); fixed {
		return nil
	}
	return n
}

// countValue is a Count's live value. Every Quantity validates its Count at card
// init, so there is no unset case to default here.
func countValue(n Count, ctx *EffectContext) int {
	return n.Value(ctx)
}

// Takes takes N cards, or as many as the zone holds when it holds fewer — the
// plain "purge 2 cards", or "for each friendly Shard, shuffle a card" when N
// scales with the board (Shard of Life).
type Takes struct {
	// N is how many cards to take; it must be set.
	N Count
}

// picks is the count's live value.
func (q Takes) picks(ctx *EffectContext) (int, bool) { return countValue(q.N, ctx), true }

// object renders "2 cards", or the Selection's own phrase for a single card.
func (q Takes) object(noun, single string) string { return countedObject(q.N, noun, single) }

// leadIn is the board count this opens the sentence with, if any.
func (q Takes) leadIn() Count { return countedLeadIn(q.N) }

// fixed is the constant behind this count, if it is one.
func (q Takes) fixed() (int, bool) { return countFixed(q.N) }

// optional is false: Takes asks for as many as the pool allows.
func (Takes) optional() bool { return false }

// validate rejects a Takes whose count was left unset or is below one.
func (q Takes) validate() error { return validateCardCount("Takes", q.N) }

// UpTo takes at most N cards, letting the controller stop early — "purge up to 2
// cards" (Creeping Oblivion). Pair it with a declinable Selection, which is what
// gives the controller the stop.
type UpTo struct {
	// N is the ceiling; it must be set. A ceiling of one reads "you may" rather
	// than "up to 1".
	N Count
}

// picks is the ceiling's live value.
func (q UpTo) picks(ctx *EffectContext) (int, bool) { return countValue(q.N, ctx), true }

// object renders "up to 2 cards". A ceiling of one adds nothing to the Selection's
// own phrase — the verb reads "you may purge a card" instead.
func (q UpTo) object(noun, single string) string {
	counted := countedObject(q.N, noun, single)
	if counted == single {
		return single
	}
	return "up to " + counted
}

// leadIn is the board count this opens the sentence with, if any.
func (q UpTo) leadIn() Count { return countedLeadIn(q.N) }

// fixed is the constant behind this ceiling, if it is one.
func (q UpTo) fixed() (int, bool) { return countFixed(q.N) }

// optional is true: the ceiling is the point of the node.
func (UpTo) optional() bool { return true }

// validate rejects an UpTo whose ceiling was left unset or is below one.
func (q UpTo) validate() error { return validateCardCount("UpTo", q.N) }

// AnyNumber takes as many matching cards as the controller likes, one at a time,
// until they decline (Not Finished with You). Pair it with a declinable Selection.
type AnyNumber struct{}

// picks is unbounded: the verb stops when a pick comes back empty.
func (AnyNumber) picks(*EffectContext) (int, bool) { return 0, false }

// object renders "any number of cards".
func (AnyNumber) object(noun, _ string) string { return "any number of " + noun + "s" }

// leadIn is none: "any number" is already the whole quantity phrase.
func (AnyNumber) leadIn() Count { return nil }

// fixed is never known: "any number" is settled only as the controller answers.
func (AnyNumber) fixed() (int, bool) { return 0, false }

// optional is true: "any number" includes none.
func (AnyNumber) optional() bool { return true }

// validate always passes: "any number" has no number to leave unset.
func (AnyNumber) validate() error { return nil }
