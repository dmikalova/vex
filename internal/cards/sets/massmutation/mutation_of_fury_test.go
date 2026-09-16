package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mutation of Fury
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Choose a creature - the chosen creature gains assault 3 and the Mutant trait until the start of your next turn.
func TestMutationOfFury(t *testing.T) {
	var beast, foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(MutationOfFury),
			InPlay: ct.Cards(
				ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(2))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(20))),
			),
		},
	})

	h.P1.Play(MutationOfFury)
	h.P1.ClickCard(beast)

	if !h.Game().HasTrait(beast.ID(), card.Traits.Mutant) {
		t.Error("the chosen creature should have gained the Mutant trait")
	}

	// The gained assault deals its 3 before the 2 fight damage: 5 on the 20-power foe.
	h.P1.Fight(beast, foe)
	h.Expect(foe).At(ct.PlayArea).Damage(5)
}

// TestMutationOfFurySpan verifies the granted Mutant trait lasts through the
// opponent's whole turn and lifts only at the start of the controller's next turn.
func TestMutationOfFurySpan(t *testing.T) {
	var beast ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(MutationOfFury),
			InPlay: ct.Cards(
				ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(2))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(20)),
			),
		},
	})

	h.P1.Play(MutationOfFury)
	h.P1.ClickCard(beast)
	if !h.Game().HasTrait(beast.ID(), card.Traits.Mutant) {
		t.Fatal("the chosen creature should have gained the Mutant trait")
	}

	h.P1.EndTurn()
	if !h.Game().HasTrait(beast.ID(), card.Traits.Mutant) {
		t.Error("the Mutant trait should survive the opponent's turn")
	}
	h.P2.EndTurn()
	if h.Game().HasTrait(beast.ID(), card.Traits.Mutant) {
		t.Error("the Mutant trait should lift at the start of the controller's next turn")
	}
}
