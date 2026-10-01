package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Neutron Shark
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Mutant
//
//	Play/Fight/Reap: Destroy an enemy creature or artifact and a friendly creature or artifact. Discard the top card of your deck. If the discarded card is not a Logos card, repeat this effect.
func TestNeutronShark(t *testing.T) {
	t.Run("stops once the discarded card is a Logos card", func(t *testing.T) {
		var friend, foe, logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(NeutronShark),
				InPlay: ct.Cards(
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Deck: ct.Cards(ct.Bind(&logos, Batdrone)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		// The enemy side has a single creature, so that destruction resolves without
		// a prompt; the friendly side offers Neutron Shark itself alongside its
		// neighbour. The discard then turns up a Logos card, ending the loop.
		h.P1.Play(NeutronShark)
		h.P1.ClickCard(friend)

		h.Expect(foe).At(ct.Discard)
		h.Expect(friend).At(ct.Discard)
		h.Expect(logos).At(ct.Discard)
	})

	// The official ruling (docs/keyforge-master-rulebook.md, "Neutron Shark"):
	// "The destruction of Neutron Shark and the first enemy creature resolves
	// fully before you proceed with the rest of Neutron Shark's ability. Once
	// Neutron Shark is no longer in play, its Play/After Fight/After Reap:
	// ability cannot be resolved again."
	t.Run("stops once Neutron Shark itself has left play", func(t *testing.T) {
		var foe1, foe2, top, second, third ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(NeutronShark),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&second, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&third, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				ct.Bind(&foe2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		// P1 controls nothing else, so Neutron Shark is the only legal friendly
		// target and that half resolves without a prompt. The started resolution
		// finishes — one card comes off the top of the deck — but the discarded
		// card being non-Logos cannot start another pass, because the source of
		// the ability is no longer in play.
		h.P1.Play(NeutronShark)
		h.P1.ClickCard(foe1)

		h.Expect(foe1).At(ct.Discard)
		h.Expect(NeutronShark).At(ct.Discard)
		h.Expect(top).At(ct.Discard)

		h.Expect(foe2).At(ct.PlayArea)
		h.Expect(second).At(ct.Deck)
		h.Expect(third).At(ct.Deck)
	})
}
