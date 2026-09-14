package engine

import "fmt"

// CardsInZone counts the cards in one of a player's zones — deck, hand, archives,
// or discard pile — so a card that scales with a zone's size (or, via CountIs,
// gates on it) reuses one count instead of a per-zone type. The house-referenced
// hand count stays CardsInHand, whose house semantics are specific to the hand.
type CardsInZone struct {
	Zone   Zone
	Player Player
}

// validate rejects a zone a card cannot name.
func (e CardsInZone) validate() error {
	switch e.Zone {
	case Deck, Hand, Archives, Discard:
		return nil
	default:
		return fmt.Errorf("CardsInZone: Zone must be Deck, Hand, Archives, or Discard")
	}
}

// cards returns the player's ids in the named zone.
func (e CardsInZone) cards(ctx *EffectContext) []LocalID {
	p := ctx.PlayerFor(e.Player)
	switch e.Zone {
	case Deck:
		return ctx.Resolver.Deck(p)
	case Hand:
		return ctx.Resolver.Hand(p)
	case Archives:
		return ctx.Resolver.Archives(p)
	default: // Discard
		return ctx.Resolver.Discard(p)
	}
}

// Value returns how many cards the chosen player has in the zone.
func (e CardsInZone) Value(ctx *EffectContext) int { return len(e.cards(ctx)) }

// possessive names whose zone it is, for the shared text.
func (e CardsInZone) possessive() string {
	if e.Player == Opponent {
		return "your opponent's"
	}
	return "your"
}

// CountText renders the singular noun the "for each" clause repeats, e.g. "card in
// your opponent's archives".
func (e CardsInZone) CountText() string {
	return "card in " + e.possessive() + " " + e.Zone.noun()
}

// CountClause renders the clause CountIs puts after "if", e.g. "you have 5 or
// fewer cards in your deck".
func (e CardsInZone) CountClause(quantity string, _ bool) string {
	has := "you have"
	if e.Player == Opponent {
		has = "your opponent has"
	}
	return has + " " + quantity + " cards in " + e.possessive() + " " + e.Zone.noun()
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
