package engine

import "strings"

// Revealing cards shows them from a hand to both players and records them in the
// log, turning hidden information public. A card reveals cards so that what
// follows can be trusted — you reveal the Mars cards you are drawing for, or an
// opponent's whole hand before discarding from it — which is why the printed text
// is careful about which cards are shown.
//
// A House narrows the reveal to cards of that house (the wording "reveal any
// number of Mars cards"): a player would only ever reveal cards that help them,
// so every matching card is revealed. An unset House reveals the whole hand.
type Reveal struct {
	Player Player
	House  House
}

// validate rejects a Reveal whose player was left unset.
func (e Reveal) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("Reveal")
	}
	return nil
}

// Text renders the effect, e.g. "reveal any number of Mars cards from your hand"
// or "reveal your opponent's hand".
func (e Reveal) Text() string {
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
func (e Reveal) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	var names []string
	for _, id := range ctx.Resolver.Hand(owner) {
		if e.House == HouseNone || ctx.Resolver.House(id) == e.House {
			names = append(names, ctx.Resolver.Name(id))
		}
	}
	ctx.Produced.Revealed = len(names)
	if len(names) > 0 {
		ctx.Resolver.Logf(
			"%s reveals %s",
			ctx.Resolver.PlayerName(owner),
			strings.Join(names, ", "),
		)
	}
}

// CardsRevealed counts the cards the most recent Reveal showed — the "for each
// card revealed this way" clause. Reveal records the tally on the context, so
// pairing it after a Reveal lets an effect scale with the reveal.
type CardsRevealed struct{}

// Value returns how many cards the preceding Reveal showed.
func (CardsRevealed) Value(ctx *EffectContext) int { return ctx.Produced.Revealed }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsRevealed) CountText() string { return "card revealed this way" }
