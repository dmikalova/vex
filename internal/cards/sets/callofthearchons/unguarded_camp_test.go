package callofthearchons_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
)

// Unguarded Camp
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each creature you have in excess of your opponent, a friendly creature captures 1 Æmber from your opponent. Each creature cannot capture more than 1 Æmber this way.
func TestUnguardedCamp(t *testing.T) {
	var camp, a, b, c ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Brobnar,
			Hand:  ct.Cards(ct.Bind(&camp, callofthearchons.UnguardedCamp)),
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				ct.Bind(&c, ct.Creature(ct.OfHouse(card.House.Brobnar))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Creature()),
			Amber:  5,
		},
	})

	// Three creatures against one is an excess of two, so two different friendly
	// creatures each capture 1 Æmber.
	h.P1.Play(camp)
	h.P1.ClickCard(a)
	h.P1.ClickCard(b)

	h.Expect(a).AmberOn(1)
	h.Expect(b).AmberOn(1)
	h.Expect(c).AmberOn(0)
	h.P2.ExpectAmber(3)
	h.P1.ExpectAmber(1) // the card's own Æmber bonus
}

// TestUnguardedCampNoExcess covers the capture not happening at all when the
// controller has no more creatures than their opponent.
func TestUnguardedCampNoExcess(t *testing.T) {
	var camp, mine ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			Hand:   ct.Cards(ct.Bind(&camp, callofthearchons.UnguardedCamp)),
			InPlay: ct.Cards(ct.Bind(&mine, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Creature(), ct.Creature()),
			Amber:  5,
		},
	})

	h.P1.Play(camp)

	h.Expect(mine).AmberOn(0)
	h.P2.ExpectAmber(5)
}
