package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Quantum Fingertrap
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Swap the positions of two creatures in a battleline.
func TestQuantumFingertrap(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			InPlay: ct.Cards(
				QuantumFingertrap,
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Logos))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Logos))),
			),
		},
	})

	h.P1.UseAction(QuantumFingertrap)
	h.P1.ClickCard(a)

	// Both creatures remain in play; using the action exhausts the artifact.
	h.Expect(a).At(ct.PlayArea)
	h.Expect(b).At(ct.PlayArea)
	h.Expect(QuantumFingertrap).At(ct.PlayArea).Exhausted()
}
