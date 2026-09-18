package engine

import "fmt"

// Count computes a number from live game state, for effects whose magnitude
// scales with the board — e.g. "for each key your opponent has forged, gain 1
// Æmber". It both yields the value (Value) and renders the leading "for each ..."
// clause of the effect's text (CountText), so the printed card and its behavior
// share one source, exactly like Effect and Condition.
type Count interface {
	Value(ctx *EffectContext) int
	CountText() string
}

// forEach front-loads a Per count onto an effect's body as a leading "for each
// ..." clause, so the sentence reads subject-first and forward (card-wording rule
// 9). With no count it returns body unchanged.
func forEach(per Count, body string) string {
	if per == nil {
		return body
	}
	return "for each " + countLeadText(per) + ", " + body
}

// scaled multiplies a base amount by a Per count when one is set — the value
// companion to forEach's text. Every effect whose magnitude is "N for each ..."
// (Draw, GainAember, DealDamage, AddPowerCounter, the economy amounts) scales the
// same way, so they share this rather than each re-testing per for nil.
func scaled(base int, per Count, ctx *EffectContext) int {
	if per == nil {
		return base
	}
	return base * per.Value(ctx)
}

// Fixed is a Count of a constant number — a repetition that always runs the same
// number of times whatever the board (RepeatedFight fights a fixed number of
// times). It lets every Times field be a Count, whether or not the count scales.
type Fixed int

// Value returns the constant, ignoring game state.
func (n Fixed) Value(*EffectContext) int { return int(n) }

// CountText is unused: a Fixed never leads a "for each" clause; the effects that
// take one print its number directly.
func (n Fixed) CountText() string { return "" }

// fixedValue reads the constant of a Count that must be a Fixed, for text built
// before resolution (no game state yet). Effects that call it validate their count
// is a Fixed, so a non-Fixed reads as zero and fails that check.
func fixedValue(c Count) int {
	f, _ := c.(Fixed)
	return int(f)
}

// leadingCounter is a Count whose "for each" clause names its subject when the
// clause leads the sentence, where a trailing "it" would read as a forward
// reference — NeighborsOfThis is "neighbor it has" trailing but "neighbor <self>
// has" leading.
type leadingCounter interface {
	leadingCountText() string
}

// countLeadText renders per's "for each" noun for a leading clause, preferring a
// leadingCounter's named form over the trailing CountText.
func countLeadText(per Count) string {
	if lc, ok := per.(leadingCounter); ok {
		return lc.leadingCountText()
	}
	return per.CountText()
}

// cardinalCounter is the optional capability of a Count that renders itself as a
// cardinal "the number of …" phrase, for a clause that compares a value against
// the count rather than repeating "for each …".
type cardinalCounter interface {
	cardinalCountText() string
}

// cardinalCountText renders a Count as "the number of …", preferring a
// cardinalCounter's own phrasing over its singular "for each" noun.
func cardinalCountText(c Count) string {
	if cc, ok := c.(cardinalCounter); ok {
		return cc.cardinalCountText()
	}
	return "the number of " + c.CountText()
}

// ForgedKeys counts the keys a player has forged — the running "for each forged
// key your opponent has" tally cards like Dr. Escotera scale by.
type ForgedKeys struct{ Player Player }

// Value returns the named player's forged-key count.
func (e ForgedKeys) Value(ctx *EffectContext) int {
	return ctx.Resolver.Keys(ctx.PlayerFor(e.Player))
}

// CountText renders the singular noun the "for each" clause repeats.
func (e ForgedKeys) CountText() string {
	if e.Player == Opponent {
		return "forged key your opponent has"
	}
	return "forged key you have"
}

// PurgedCards counts every card set aside in the purge pile across both players —
// the running "+1 power for each purged card" tally Noname scales its power by.
type PurgedCards struct{}

// Value returns the combined size of both players' purge piles.
func (PurgedCards) Value(ctx *EffectContext) int {
	return len(ctx.Resolver.Purge(ctx.Controller)) + len(ctx.Resolver.Purge(ctx.Opponent()))
}

// CountText renders the singular noun the "for each" clause repeats.
func (PurgedCards) CountText() string { return "purged card" }

// HousesInPlay counts the distinct houses represented among all cards in play,
// optionally excluding one house — Free Markets pays out per house other than
// its own Sanctum.
type HousesInPlay struct{ Except House }

// Value counts the distinct houses on the board, skipping the excepted one.
func (e HousesInPlay) Value(ctx *EffectContext) int {
	seen := map[House]bool{}
	for _, p := range [2]int{0, 1} {
		for _, id := range append(ctx.Resolver.Battleline(p), ctx.Resolver.Artifacts(p)...) {
			if h := ctx.Resolver.House(id); h != e.Except {
				seen[h] = true
			}
		}
	}
	return len(seen)
}

// CountText renders the singular noun the "for each" clause repeats.
func (e HousesInPlay) CountText() string {
	if e.Except != HouseNone {
		return "house represented among cards in play, except for " + e.Except.String()
	}
	return "house represented among cards in play"
}

// HousesAmong counts the distinct houses represented among a chosen set of
// in-play cards — a player's creatures, all creatures, or all cards in play — for
// effects that scale with how many houses share the board (Trust No One,
// Quadracorder, Galactic Census, Forging an Alliance). It is the player- and
// type-scoped companion to HousesInPlay, which counts every card in play minus one
// excepted house (Free Markets); a houseless card counts toward no house.
type HousesAmong struct {
	// Player names whose cards to survey: Controller (friendly), Opponent (enemy),
	// or EachPlayer (both).
	Player Player
	// Type filters the surveyed cards; the zero value surveys any type, Creature
	// only creatures.
	Type CardType
}

// Value counts the distinct houses on the surveyed cards.
func (e HousesAmong) Value(ctx *EffectContext) int {
	var seen [NumHouses]bool
	n := 0
	for _, p := range e.players(ctx) {
		for _, id := range e.set(ctx, p) {
			if h := ctx.Resolver.House(id); h != HouseNone && !seen[h] {
				seen[h] = true
				n++
			}
		}
	}
	return n
}

// players returns the sides to survey: both under EachPlayer, else the named one.
func (e HousesAmong) players(ctx *EffectContext) []int {
	if e.Player == EachPlayer {
		return []int{0, 1}
	}
	return []int{ctx.PlayerFor(e.Player)}
}

// set returns a player's in-play ids the type filter surveys: the battleline for
// creatures, or both rows when the type is unset.
func (e HousesAmong) set(ctx *EffectContext, p int) []LocalID {
	if e.Type == Creature {
		return ctx.Resolver.Battleline(p)
	}
	return append(ctx.Resolver.Battleline(p), ctx.Resolver.Artifacts(p)...)
}

// scope names the surveyed set as a plural noun the text roles share: "friendly
// creatures", "enemy creatures", "creatures in play", or "cards in play".
func (e HousesAmong) scope() string {
	noun := "cards"
	if e.Type == Creature {
		noun = "creatures"
	}
	switch e.Player {
	case Controller:
		return "friendly " + noun
	case Opponent:
		return "enemy " + noun
	default:
		return noun + " in play"
	}
}

// CountText renders the singular noun the "for each" clause repeats. Any Max cap
// is silent, so it does not appear in the text.
func (e HousesAmong) CountText() string {
	return "house represented among " + e.scope()
}

// CountClause renders the clause CountIs puts after "if", e.g. "3 or more houses
// are represented among friendly creatures". It names the same surveyed set as
// CountText, so the "for each" and "if" voices cannot drift apart.
func (e HousesAmong) CountClause(quantity string, plural bool) string {
	noun, verb := "house", "is"
	if plural {
		noun, verb = "houses", "are"
	}
	return fmt.Sprintf("%s %s %s represented among %s", quantity, noun, verb, e.scope())
}

// UnforgedKeys counts the keys a player still has to forge — the measure of how
// far they are from winning, which Mushroom Man grows on.
type UnforgedKeys struct{ Player Player }

// Value returns how many of the player's keys are still unforged.
func (e UnforgedKeys) Value(ctx *EffectContext) int {
	return KeysToWin - ctx.Resolver.Keys(ctx.PlayerFor(e.Player))
}

// CountText renders the singular noun the "for each" clause repeats.
func (e UnforgedKeys) CountText() string {
	if e.Player == Opponent {
		return "unforged key your opponent has"
	}
	return "unforged key you have"
}

// DamageOnThis counts the damage sitting on the source card, so a card can scale
// with how hurt it is (Angwish charges its opponent +1 key cost per damage on it).
type DamageOnThis struct{}

// Value returns the damage on the source card.
func (DamageOnThis) Value(ctx *EffectContext) int { return ctx.damageOn(ctx.Source) }

// CountText renders the singular noun the "for each" clause repeats.
func (DamageOnThis) CountText() string { return "damage on it" }

// portionPhraser is the optional capability of a Loss that can also phrase a
// fraction of a count rather than a pool — "half its power, rounded down". Fraction
// implements it, so the same portion vocabulary serves both LoseAember{By:
// HalfRoundedDown} and PowerOfChosen{Of: HalfRoundedDown}.
type portionPhraser interface {
	countPhrase(noun string) string
}

// PowerOfChosen is the power of the creature in context (ctx.It) — the creature a
// fight or a ChooseCreatureThen put in context. Mindworm makes the creature it
// fights deal damage equal to its power to each of its neighbors. Of takes a
// fraction of that power (Of: HalfRoundedDown is half, rounded down — The Flex);
// the zero value is the full power.
type PowerOfChosen struct {
	Of Loss
}

// Value returns the context creature's power, reduced by the Of fraction when set,
// or zero when no creature is in context.
func (e PowerOfChosen) Value(ctx *EffectContext) int {
	if !ctx.HasIt {
		return 0
	}
	power := ctx.powerOf(ctx.It)
	if e.Of != nil {
		power = e.Of.lose(power)
	}
	return power
}

// CountText renders the amount as it reads in an "equal to" clause, e.g. "its
// power" or "half its power, rounded down".
func (e PowerOfChosen) CountText() string {
	if p, ok := e.Of.(portionPhraser); ok {
		return p.countPhrase("its power")
	}
	return "its power"
}

// TraitsOfChosen counts the traits of the creature in context (ctx.It) — the
// creature a ChooseCreatureThen just picked. Entropic Swirl acts once per trait
// the chosen creature has.
type TraitsOfChosen struct{}

// Value returns the number of traits on the context creature, or zero when no
// creature is in context.
func (TraitsOfChosen) Value(ctx *EffectContext) int {
	if !ctx.HasIt {
		return 0
	}
	return ctx.Resolver.TraitCount(ctx.It)
}

// CountText renders the singular noun the "for each" clause repeats.
func (TraitsOfChosen) CountText() string { return "trait that creature has" }

// BonusIconsOfChosen counts the bonus icons on the card in context (ctx.It) — the
// card an effect just discarded or revealed. Mindfire steals 1 Æmber for each
// bonus icon on the card it discarded. Noun names the card in the text so the
// clause reads "the discarded card" rather than a bare "it".
type BonusIconsOfChosen struct {
	Noun ItNoun
}

// Value returns the number of bonus icons on the context card, or zero when no
// card is in context.
func (e BonusIconsOfChosen) Value(ctx *EffectContext) int {
	if !ctx.HasIt {
		return 0
	}
	return ctx.Resolver.BonusIconCount(ctx.It)
}

// CountText renders the singular noun the "for each" clause repeats, e.g. "bonus
// icon on the discarded card".
func (e BonusIconsOfChosen) CountText() string {
	return "bonus icon on " + e.Noun.noun()
}

// CopiesInDiscard counts the cards in the controller's discard pile sharing the
// source card's name — a card that pays off for having been played before
// (Routine Job). The card being resolved is not in the discard pile yet, so it
// never counts itself.
type CopiesInDiscard struct{}

// Value counts the copies of the source card in the controller's discard pile.
func (CopiesInDiscard) Value(ctx *EffectContext) int {
	name := ctx.Resolver.Name(ctx.Source)
	n := 0
	for _, id := range ctx.Resolver.Discard(ctx.Controller) {
		if ctx.Resolver.Name(id) == name {
			n++
		}
	}
	return n
}

// CountText renders the singular noun the "for each" clause repeats.
func (CopiesInDiscard) CountText() string {
	return "copy of " + SelfName + " in your discard pile"
}

// TurnCount counts one of the engine's turn-history tallies for a player — the
// creatures they played on their previous turn (Lifeweb), the enemy creatures
// destroyed in a fight this turn (The Warchest). One Count over a TurnStat rather
// than a node per tally, so asking a new question about a turn is a new enum
// value and its noun.
type TurnCount struct {
	Player Player
	Of     TurnStat
}

// Value reads the tally for the player the count names.
func (c TurnCount) Value(ctx *EffectContext) int {
	return ctx.Resolver.TurnHistory(ctx.PlayerFor(c.Player), c.Of)
}

// CountText renders the singular noun the "for each" clause repeats.
func (c TurnCount) CountText() string { return turnStatNoun[c.Of] }

// CountClause renders the "if ..." clause CountIs needs, e.g. "if your opponent
// played 3 or more creatures on their previous turn" or, for the enemy-destroyed
// tally, "if 3 or more enemy creatures have been destroyed this turn".
func (c TurnCount) CountClause(quantity string, plural bool) string {
	switch c.Of {
	case EnemyCreaturesDestroyed:
		noun := "enemy creature has"
		if plural {
			noun = "enemy creatures have"
		}
		return fmt.Sprintf("%s %s been destroyed this turn", quantity, noun)
	default:
		subject, possessive := "you played", "your"
		if c.Player == Opponent {
			subject, possessive = "your opponent played", "their"
		}
		noun := "creature"
		if plural {
			noun = "creatures"
		}
		return fmt.Sprintf("%s %s %s on %s previous turn", subject, quantity, noun, possessive)
	}
}

// NeighborsOfThis counts the battleline neighbors of the creature holding the
// ability — 0, 1, or 2. Knoxx grows by 3 power for each neighbor it has. Its
// value only changes when the battleline does, so the destroyed sweep already
// runs each time it can shift.
type NeighborsOfThis struct{}

// Value counts the source creature's immediate battleline neighbors.
func (c NeighborsOfThis) Value(ctx *EffectContext) int {
	return len(neighbors(ctx, ctx.Source))
}

// CountText renders the singular noun the "for each" clause repeats.
func (c NeighborsOfThis) CountText() string {
	return "neighbor it has"
}

// leadingCountText names the source creature when the clause leads the sentence,
// where a trailing "it" would be a forward reference — Nyzyk Resonator.
func (c NeighborsOfThis) leadingCountText() string {
	return "neighbor " + SelfName + " has"
}

// CombinedPowerOfNeighborsWithout sums the current power of the source creature's
// battleline neighbors that do not carry Without — Picaroon's X is the combined
// power of its non-Changeling neighbors. It reads power live, so a neighbor's buff
// or a change to the battleline shifts the total.
type CombinedPowerOfNeighborsWithout struct {
	Without Trait
}

// Value sums the power of the source's neighbors that lack the excluded trait.
func (c CombinedPowerOfNeighborsWithout) Value(ctx *EffectContext) int {
	sum := 0
	for _, id := range neighbors(ctx, ctx.Source) {
		if !ctx.Resolver.HasTrait(id, c.Without) {
			sum += ctx.Resolver.Power(id)
		}
	}
	return sum
}

// CountText renders the singular noun a "for each" clause would repeat.
func (c CombinedPowerOfNeighborsWithout) CountText() string {
	return "combined power of " + SelfName + "'s non-" + c.Without.String() + " neighbors"
}

// cardinalCountText renders the value as a standalone phrase, for the "X is …"
// power line rather than a "for each …" clause.
func (c CombinedPowerOfNeighborsWithout) cardinalCountText() string {
	return "the combined power of " + SelfName + "'s non-" + c.Without.String() + " neighbors"
}

// NeighborsMatching counts the battleline neighbors of the creature in context
// (ctx.It) that House admits — Thorium Plasmate deals 2 damage to a moved creature
// for each neighbor of that card's house (Houses.Contextual).
type NeighborsMatching struct {
	House HouseMatcher
}

// Value counts the context creature's immediate neighbors House admits.
func (c NeighborsMatching) Value(ctx *EffectContext) int {
	if !ctx.HasIt {
		return 0
	}
	n := 0
	for _, id := range neighbors(ctx, ctx.It) {
		if c.House.matches(ctx, id) {
			n++
		}
	}
	return n
}

// CountText renders the singular noun the "for each" clause repeats.
func (c NeighborsMatching) CountText() string {
	return c.House.qualify("neighbor")
}
