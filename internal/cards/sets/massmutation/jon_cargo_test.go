package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// J.O.N. Cargo
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	Reap: Discard the top card of your deck. Reveal your hand. Archive each card of that card's house from your hand.
func TestJONCargo(t *testing.T) {
	var jon, top, handMatch1, handMatch2, handOther ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(ct.Bind(&jon, JONCargo)),
			Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			Hand: ct.Cards(
				ct.Bind(&handMatch1, ct.Creature(ct.OfHouse(card.House.Logos))),
				ct.Bind(&handMatch2, ct.Tactic(ct.OfHouse(card.House.Logos))),
				ct.Bind(&handOther, ct.Creature(ct.OfHouse(card.House.Brobnar))),
			),
		},
	})

	h.P1.Reap(jon)

	// The discarded top card fixes the house; hand cards of that house are archived.
	h.Expect(top).At(ct.Discard)
	h.Expect(handMatch1).At(ct.Archives)
	h.Expect(handMatch2).At(ct.Archives)
	h.Expect(handOther).At(ct.Hand)
}
