package engine

import (
	"fmt"
	"strings"
)

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
	return "for each " + per.CountText() + ", " + body
}

// OpponentForgedKeys counts the keys the controller's opponent has forged.
type OpponentForgedKeys struct{}

// Value returns the opponent's forged-key count.
func (OpponentForgedKeys) Value(ctx *EffectContext) int {
	return ctx.Resolver.Keys(ctx.Opponent())
}

// CountText renders the singular noun the "for each" clause repeats.
func (OpponentForgedKeys) CountText() string { return "forged key your opponent has" }

// ExcessCreatures counts how many more creatures one player controls than the
// other (never below zero). Player names whose excess is counted: Opponent for
// "each creature your opponent controls in excess of you" (Glorious Few),
// Controller for "each creature you have in excess of your opponent"
// (Unguarded Camp).
type ExcessCreatures struct{ Player Player }

// Value returns the named player's creature count minus the other's, floored at 0.
func (e ExcessCreatures) Value(ctx *EffectContext) int {
	more := ctx.PlayerFor(e.Player)
	return max(0, len(ctx.Resolver.Battleline(more))-len(ctx.Resolver.Battleline(1-more)))
}

// CountText renders the singular noun the "for each" clause repeats.
func (e ExcessCreatures) CountText() string {
	if e.Player == Opponent {
		return "creature your opponent controls in excess of you"
	}
	return "creature you have in excess of your opponent"
}

// CardsDestroyed counts the cards the most recent destruction in this resolution
// removed from play — the "for each card destroyed this way" tally (Oath of
// Poverty gains 2 Æmber for each artifact it destroyed).
type CardsDestroyed struct{}

// Value returns how many cards the preceding Destroy actually removed.
func (CardsDestroyed) Value(ctx *EffectContext) int { return ctx.Produced.TotalDestroyed() }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsDestroyed) CountText() string { return "card destroyed this way" }

// CardsInArchives counts the cards in a player's archives.
type CardsInArchives struct{ Player Player }

// Value returns how many cards the chosen player has archived.
func (e CardsInArchives) Value(ctx *EffectContext) int {
	return len(ctx.Resolver.Archives(ctx.PlayerFor(e.Player)))
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CardsInArchives) CountText() string {
	if e.Player == Opponent {
		return "card in your opponent's archives"
	}
	return "card in your archives"
}

// CardsInHand counts the cards in a player's hand of a referenced house — the
// house chosen this turn, the active house, or the house of the card in context
// (A Fair Game counts the opponent's hand cards sharing its just-discarded deck
// card's house). A referenced house that resolves to none counts zero.
type CardsInHand struct {
	Player Player
	House  HouseChoice
}

// Value counts the player's hand cards matching the referenced house.
func (e CardsInHand) Value(ctx *EffectContext) int {
	hand := ctx.Resolver.Hand(ctx.PlayerFor(e.Player))
	if e.House == AnyHouse {
		return len(hand)
	}
	house := e.House.resolveHouse(ctx)
	if house == HouseNone {
		return 0
	}
	n := 0
	for _, id := range hand {
		if ctx.Resolver.House(id) == house {
			n++
		}
	}
	return n
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CardsInHand) CountText() string {
	switch e.House {
	case AnyHouse:
		if e.Player == Opponent {
			return "card in your opponent's hand"
		}
		return "card in your hand"
	case TheContextualHouse:
		return "card of the discarded card's house revealed this way"
	case TheActiveHouse:
		return "card of the active house in their hand"
	default:
		return "card of the chosen house in their hand"
	}
}

// InPlay selects the cards a player has in play that match its filters — of a
// given type and house — and serves two roles from one description. As a Count (a
// Per clause) it yields the match count and renders the repeated "for each ..."
// noun; as a Condition it is met when the count reaches Amount, which defaults to
// one. This unifies the several friendly-creature counts and conditions, e.g.
// InPlay{Player: Controller, Type: Creature} or the house-filtered
// InPlay{Player: Controller, Type: Creature, House: Mars}.
type InPlay struct {
	Player Player
	// Type filters by card type; the zero value counts any type.
	Type CardType
	// House filters by house; the zero value (HouseNone) counts any house.
	House House
	// Ready counts only cards that are ready (not exhausted).
	Ready bool
	// Damaged counts only creatures that have damage on them.
	Damaged bool
	// Other leaves the source card itself out of the count, which is how a card
	// counts its companions (Phylyx the Disintegrator).
	Other bool
	// Name filters by printed card name, and replaces the rendered noun with it:
	// "if there are no Ancient Bears in play" (Bear Flute).
	Name string
	// None inverts the Condition role: it is met when nothing matches. It reads as
	// its own word rather than as Amount 0 because "if there are no X" is a
	// different sentence from "if there are N X", not the N = 0 case of it.
	None bool
	// Amount is the minimum the Condition role requires; zero means at least one.
	Amount int
}

// Value counts the matching cards the player has in play.
func (e InPlay) Value(ctx *EffectContext) int {
	n := 0
	for _, id := range e.set(ctx) {
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			continue
		}
		if e.Ready && ctx.Resolver.Exhausted(id) {
			continue
		}
		if e.Damaged && ctx.Resolver.Damage(id) == 0 {
			continue
		}
		if e.Other && id == ctx.Source {
			continue
		}
		if e.Name != "" && ctx.Resolver.Name(id) != e.Name {
			continue
		}
		n++
	}
	return n
}

// Met reports whether at least Amount (default one) matching cards are in play,
// or, under None, that none are.
func (e InPlay) Met(ctx *EffectContext) bool {
	if e.None {
		return e.Value(ctx) == 0
	}
	return e.Value(ctx) >= e.threshold()
}

// set returns the player's in-play ids the type filter considers: the battleline
// for creatures, the artifact row for artifacts, or both when the type is unset.
func (e InPlay) set(ctx *EffectContext) []LocalID {
	if e.Player == EachPlayer {
		return append(e.playerSet(ctx, 0), e.playerSet(ctx, 1)...)
	}
	return e.playerSet(ctx, ctx.PlayerFor(e.Player))
}
func (e InPlay) playerSet(ctx *EffectContext, p int) []LocalID {
	switch e.Type {
	case Creature:
		return ctx.Resolver.Battleline(p)
	case Artifact:
		return ctx.Resolver.Artifacts(p)
	default:
		return append(ctx.Resolver.Battleline(p), ctx.Resolver.Artifacts(p)...)
	}
}

// threshold is the Condition's required count, defaulting to one.
func (e InPlay) threshold() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// who renders the controlling side as "friendly" or "enemy".
func (e InPlay) who() string {
	switch e.Player {
	case Opponent:
		return "enemy"
	case EachPlayer:
		return ""
	default:
		return "friendly"
	}
}

// typeNoun renders the filtered type as a noun.
func (e InPlay) typeNoun() string {
	switch e.Type {
	case Creature:
		return "creature"
	case Artifact:
		return "artifact"
	default:
		return "card"
	}
}

// noun renders the "<side> [ready ][house ]<type>" phrase the text roles share.
func (e InPlay) noun() string {
	if e.Name != "" {
		return e.Name
	}
	parts := []string{}
	if e.Other {
		parts = append(parts, "other")
	}
	if who := e.who(); who != "" {
		parts = append(parts, who)
	}
	if e.Ready {
		parts = append(parts, "ready")
	}
	if e.Damaged {
		parts = append(parts, "damaged")
	}
	if e.House != HouseNone {
		parts = append(parts, e.House.String())
	}
	return strings.Join(append(parts, e.typeNoun()), " ")
}

// CountText renders the singular noun the "for each" clause repeats. A
// house-filtered count reads "friendly Mars creature"; an unfiltered one adds "in
// play" to distinguish it from cards in hand.
func (e InPlay) CountText() string {
	if e.House != HouseNone && e.Player != EachPlayer {
		return e.noun()
	}
	return e.noun() + " in play"
}

// CondText renders the condition, e.g. "if there is a friendly creature in play"
// or "if there are 2 friendly creatures in play".
func (e InPlay) CondText() string {
	if e.None {
		return fmt.Sprintf("if there are no %s in play", plural(0, e.noun()))
	}
	if n := e.threshold(); n > 1 {
		return fmt.Sprintf("if there are %d %s in play", n, plural(n, e.noun()))
	}
	return fmt.Sprintf("if there is %s in play", indefinite(e.noun()))
}

// CardsPlayed counts the cards of a house a player has played this turn. Like
// InPlay it serves two roles from one description: as a Count it yields the tally
// and renders a "for each ... you have played this turn" noun; as a Condition it is
// met once Amount have been played (Epic Quest fires after seven), defaulting to
// one. This replaces a bespoke "played at least N of a house" condition.
type CardsPlayed struct {
	Player Player
	House  House
	// Amount is the minimum the Condition role requires; zero means at least one.
	Amount int
}

// Value counts the player's cards of the house played this turn.
func (e CardsPlayed) Value(ctx *EffectContext) int {
	return countOfHouse(ctx, ctx.Resolver.PlayedThisTurn(ctx.PlayerFor(e.Player)), e.House)
}

// countOfHouse counts how many of the ids belong to a house — the filter a
// turn-log Count applies to the unfiltered record the engine keeps. An unset
// house counts every card, so a Count can ask "how many cards" as well as "how
// many Sanctum cards".
func countOfHouse(ctx *EffectContext, ids []LocalID, house House) int {
	if house == HouseNone {
		return len(ids)
	}
	n := 0
	for _, id := range ids {
		if ctx.Resolver.House(id) == house {
			n++
		}
	}
	return n
}

// Met reports whether at least Amount (default one) matching cards were played.
func (e CardsPlayed) Met(ctx *EffectContext) bool { return e.Value(ctx) >= e.threshold() }

// threshold is the Condition's required count, defaulting to one.
func (e CardsPlayed) threshold() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CardsPlayed) CountText() string {
	return e.cardNoun() + " you have played this turn"
}

// cardNoun is the noun the count repeats, house-qualified when the count filters
// by house and a plain "card" when it counts every card played.
func (e CardsPlayed) cardNoun() string {
	if e.House == HouseNone {
		return "card"
	}
	return e.House.String() + " card"
}

// CondText renders the condition, e.g. "if you have played 7 or more Sanctum cards
// this turn".
func (e CardsPlayed) CondText() string {
	return fmt.Sprintf("if you have played %d or more %ss this turn", e.threshold(), e.cardNoun())
}

// CountClause renders the clause CountIs puts after "if", e.g. "you played
// exactly 1 card this turn".
func (e CardsPlayed) CountClause(quantity string, plural bool) string {
	noun := e.cardNoun()
	if plural {
		noun += "s"
	}
	return fmt.Sprintf("you played %s %s this turn", quantity, noun)
}

// CreaturesUsed counts the creatures a player has used this turn — reaped,
// fought, or fired an "Action:" with. It reads the per-creature use tally the
// engine already keeps, so only creatures still in play are counted.
type CreaturesUsed struct {
	Player Player
}

// Value counts the player's creatures in play that have been used this turn.
func (e CreaturesUsed) Value(ctx *EffectContext) int {
	n := 0
	for _, id := range ctx.Resolver.Battleline(ctx.PlayerFor(e.Player)) {
		if ctx.Resolver.TimesUsedThisTurn(id) > 0 {
			n++
		}
	}
	return n
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CreaturesUsed) CountText() string { return "creature you used this turn" }

// CountClause renders the clause CountIs puts after "if", e.g. "you used 3 or
// more creatures this turn".
func (e CreaturesUsed) CountClause(quantity string, plural bool) string {
	noun := "creature"
	if plural {
		noun += "s"
	}
	return fmt.Sprintf("you used %s %s this turn", quantity, noun)
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

// AemberOnThis counts the Æmber sitting on the source card, so a card can grow
// with what it captures (Yxili Marauder).
type AemberOnThis struct{}

// Value returns the Æmber on the source card.
func (AemberOnThis) Value(ctx *EffectContext) int { return ctx.Resolver.AmberOn(ctx.Source) }

// CountText renders the singular noun the "for each" clause repeats.
func (AemberOnThis) CountText() string { return "Æmber on it" }

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
// played 3 or more creatures on their previous turn".
func (c TurnCount) CountClause(quantity string, plural bool) string {
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

// The two "... this way" counts below read one player's share of a tally rather
// than the whole of it, and Player names whose. Under a GainAember{Player:
// EachPlayer} — the only place they are used — the context flips to each player
// in turn, so Controller means each player counting their own losses.

// CreaturesDestroyedThisWay counts the creatures Player controlled that an
// earlier effect in this resolution destroyed. Use CardsDestroyed for the whole
// tally, both sides together.
type CreaturesDestroyedThisWay struct{ Player Player }

// Value reads that player's share of the destruction tally.
func (c CreaturesDestroyedThisWay) Value(ctx *EffectContext) int {
	return ctx.Produced.Destroyed[ctx.PlayerFor(c.Player)]
}

// CountText renders the singular noun the "for each" clause repeats.
func (c CreaturesDestroyedThisWay) CountText() string {
	who := "they"
	if c.Player == Opponent {
		who = "your opponent"
	}
	return "creature " + who + " controlled that was destroyed this way"
}

// CreaturesShuffledIntoDeckThisWay counts the creatures Player controlled that
// an earlier effect in this resolution put back into a deck — Mating Season pays
// each player for the creatures that went home.
type CreaturesShuffledIntoDeckThisWay struct{ Player Player }

// Value reads that player's share of the put-from-play tally.
func (c CreaturesShuffledIntoDeckThisWay) Value(ctx *EffectContext) int {
	return ctx.Produced.Moved[ctx.PlayerFor(c.Player)]
}

// CountText renders the singular noun the "for each" clause repeats.
func (c CreaturesShuffledIntoDeckThisWay) CountText() string {
	whose := "their"
	if c.Player == Opponent {
		whose = "your opponent's"
	}
	return "creature shuffled into " + whose + " deck this way"
}

// AemberLostThisWay counts the Æmber an earlier LoseAember in this resolution
// took from Player's pool — Shatter Storm empties your pool and then drains your
// opponent for triple what left it.
type AemberLostThisWay struct{ Player Player }

// Value reads that player's share of the Æmber-lost tally.
func (c AemberLostThisWay) Value(ctx *EffectContext) int {
	return ctx.Produced.AemberLost[ctx.PlayerFor(c.Player)]
}

// CountText renders the singular noun the "for each" clause repeats.
func (c AemberLostThisWay) CountText() string {
	who := "you"
	if c.Player == Opponent {
		who = "your opponent"
	}
	return "Æmber " + who + " lost this way"
}

// CardsReturnedThisWay counts the cards an earlier PutFromDiscard in this
// resolution recovered from the discard pile — Ortannu the Chained deals a hit
// for each Binding it returned.
type CardsReturnedThisWay struct{}

// Value reads how many cards the preceding PutFromDiscard returned.
func (CardsReturnedThisWay) Value(ctx *EffectContext) int { return ctx.Produced.Returned }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsReturnedThisWay) CountText() string { return "card returned this way" }
