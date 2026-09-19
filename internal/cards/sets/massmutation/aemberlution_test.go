package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Aemberlution
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Omega.
//	Play: Destroy each creature. Each player reveals their hand and puts each creature from their hand into play ready.
func TestAemberlution(t *testing.T) {
	t.Run(
		"destroys every creature, then both players hatch their hand creatures ready",
		func(t *testing.T) {
			var myHandCreature, foeHandCreature, ally, foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Untamed,
					Hand: ct.Cards(
						Aemberlution,
						ct.Bind(&myHandCreature, ct.Creature(ct.Power(3))),
					),
					// Omega ends the turn, so a draw step follows; give P1 a stocked
					// deck so the refill draws from the deck rather than recycling the
					// discard.
					Deck: ct.Cards(
						ct.Creature(), ct.Creature(), ct.Creature(),
						ct.Creature(), ct.Creature(), ct.Creature(),
					),
					InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(4)))),
				},
				P2: ct.Side{
					Hand:   ct.Cards(ct.Bind(&foeHandCreature, ct.Creature(ct.Power(5)))),
					InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6)))),
				},
			})

			h.P1.Play(Aemberlution)
			h.P1.ClickOption("P1 first")

			// Every creature in play when it resolved is destroyed.
			h.Expect(ally).At(ct.Discard)
			h.Expect(foe).At(ct.Discard)
			// Each player's hand creature enters their own play area ready.
			h.Expect(myHandCreature).At(ct.PlayArea).Ready()
			h.Expect(foeHandCreature).At(ct.PlayArea).Ready()
			h.Expect(Aemberlution).At(ct.Discard)
		},
	)
}
