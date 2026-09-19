package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Essence Scale
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Choose a friendly creature. Destroy the chosen creature. Ready and use a friendly creature of that card's house.
func TestEssenceScale(t *testing.T) {
	t.Run(
		"destroys a friendly creature, then readies and uses another of its house",
		func(t *testing.T) {
			var toDestroy, disAlly, logosAlly ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					InPlay: ct.Cards(
						EssenceScale,
						ct.Bind(&toDestroy, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(4))),
						ct.Bind(&disAlly, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(4))),
						// A friendly creature of a different house cannot be the one
						// readied and used — it does not share the destroyed creature's house.
						ct.Bind(&logosAlly, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(4))),
					),
				},
				P2: ct.Side{},
			})

			// Exhaust the Dis ally first, so its readying by Essence Scale is observable.
			h.P1.Reap(disAlly)
			h.P1.ExpectAmber(1)

			h.P1.UseAction(EssenceScale)
			h.P1.ClickCard(toDestroy) // destroy this friendly creature

			// The destroyed creature was Dis, so the only friendly creature that shares
			// its house — disAlly — is readied and used (it reaps again). The Logos ally
			// is never a candidate.
			h.Expect(toDestroy).At(ct.Discard)
			h.P1.ExpectAmber(2) // the second reap
			h.Expect(disAlly).Exhausted()
			h.Expect(logosAlly).Ready()
		},
	)
}
