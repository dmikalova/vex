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
	return "for each " + per.CountText() + ", " + body
}

// OpponentForgedKeys counts the keys the controller's opponent has forged.
type OpponentForgedKeys struct{}

// Value returns the opponent's forged-key count.
func (OpponentForgedKeys) Value(ctx *EffectContext) int {
	return ctx.Resolver.Keys(ctx.Opponent())
}

// CountText renders the singular noun the "for each" clause repeats.
func (OpponentForgedKeys) CountText() string { return "key your opponent has forged" }

// OpponentExcessCreatures counts how many more creatures the controller's opponent
// controls than the controller does (never below zero) — Glorious Few.
type OpponentExcessCreatures struct{}

// Value returns the opponent's creature count minus the controller's, floored at 0.
func (OpponentExcessCreatures) Value(ctx *EffectContext) int {
	return max(
		0,
		len(ctx.Resolver.Battleline(ctx.Opponent()))-len(ctx.Resolver.Battleline(ctx.Controller)),
	)
}

// CountText renders the singular noun the "for each" clause repeats.
func (OpponentExcessCreatures) CountText() string {
	return "creature your opponent controls in excess of you"
}

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
	house := e.House.resolveHouse(ctx)
	if house == HouseNone {
		return 0
	}
	n := 0
	for _, id := range ctx.Resolver.Hand(ctx.PlayerFor(e.Player)) {
		if ctx.Resolver.House(id) == house {
			n++
		}
	}
	return n
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CardsInHand) CountText() string {
	switch e.House {
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
	// Amount is the minimum the Condition role requires; zero means at least one.
	Amount int
}

// Value counts the matching cards the player has in play.
func (e InPlay) Value(ctx *EffectContext) int {
	n := 0
	for _, id := range e.set(ctx) {
		if e.House == HouseNone || ctx.Resolver.House(id) == e.House {
			n++
		}
	}
	return n
}

// Met reports whether at least Amount (default one) matching cards are in play.
func (e InPlay) Met(ctx *EffectContext) bool {
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

// noun renders the "<side> [house ]<type>" phrase the text roles share.
func (e InPlay) noun() string {
	who := e.who()
	if e.House != HouseNone {
		if who == "" {
			return fmt.Sprintf("%s %s", e.House, e.typeNoun())
		}
		return fmt.Sprintf("%s %s %s", who, e.House, e.typeNoun())
	}
	if who == "" {
		return e.typeNoun()
	}
	return fmt.Sprintf("%s %s", who, e.typeNoun())
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
	if n := e.threshold(); n > 1 {
		return fmt.Sprintf("if there are %d %ss in play", n, e.noun())
	}
	return fmt.Sprintf("if there is a %s in play", e.noun())
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
	return ctx.Resolver.CardsPlayedOfHouseThisTurn(ctx.PlayerFor(e.Player), e.House)
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
	return fmt.Sprintf("%s card you have played this turn", e.House)
}

// CondText renders the condition, e.g. "if you have played 7 or more Sanctum cards
// this turn".
func (e CardsPlayed) CondText() string {
	return fmt.Sprintf("if you have played %d or more %s cards this turn", e.threshold(), e.House)
}
