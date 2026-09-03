package engine

// Revealing cards shows them from a hand to both players and records them in the
// log, turning hidden information public. A card reveals cards so that what
// follows can be trusted — you reveal the Mars cards you are drawing for, or an
// opponent's whole hand before discarding from it — which is why the printed text
// is careful about which cards are shown.
//
// A House narrows the reveal to cards of that house (the wording "reveal any
// number of Mars cards"): the player picks which of them to show, one at a time,
// until they are done — "any number" includes none. An unset House reveals the
// whole hand, which is not a choice.
type RevealHand struct {
	Player Player
	House  House
}

// validate rejects a Reveal whose player was left unset.
func (e RevealHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("Reveal")
	}
	return nil
}

// Text renders the effect, e.g. "reveal any number of Mars cards from your hand"
// or "reveal your opponent's hand".
func (e RevealHand) Text() string {
	whose := "your"
	if e.Player == Opponent {
		whose = "your opponent's"
	}
	if e.House == HouseNone {
		return "reveal " + whose + " hand"
	}
	return "reveal any number of " + e.House.String() + " cards from " + whose + " hand"
}

// Resolve shows the matching cards, logs them, and records how many were revealed.
func (e RevealHand) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	revealed := e.reveal(ctx, owner)
	ctx.Produced.Revealed = len(revealed)
	if len(revealed) > 0 {
		ctx.Resolver.Record(CardsRevealedToAll{Player: owner, Cards: revealed})
	}
}

// reveal returns the cards actually shown: the whole hand for an unrestricted
// reveal, or the subset the controller picks out of the matching cards when the
// reveal is "any number of <house> cards".
func (e RevealHand) reveal(ctx *EffectContext, owner int) []LocalID {
	hand := ctx.Resolver.Hand(owner)
	if e.House == HouseNone {
		return hand
	}
	var remaining []LocalID
	for _, id := range hand {
		if ctx.Resolver.House(id) == e.House {
			remaining = append(remaining, id)
		}
	}
	var shown []LocalID
	for len(remaining) > 0 {
		chosen, ok := ctx.ChooseCardOptional("Choose a card to reveal", remaining)
		if !ok {
			break
		}
		shown = append(shown, chosen)
		remaining = withoutID(remaining, chosen)
	}
	return shown
}

// CardsRevealed counts the cards the most recent Reveal showed — the "for each
// card revealed this way" clause. Reveal records the tally on the context, so
// pairing it after a Reveal lets an effect scale with the reveal.
type CardsRevealed struct{}

// Value returns how many cards the preceding Reveal showed.
func (CardsRevealed) Value(ctx *EffectContext) int { return ctx.Produced.Revealed }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsRevealed) CountText() string { return "card revealed this way" }
