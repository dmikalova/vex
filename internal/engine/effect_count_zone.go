package engine

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
