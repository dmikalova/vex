package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mutation of Cunning
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Choose a creature. The chosen creature gains elusive and the Mutant trait until the start of your next turn.
func TestMutationOfCunning(t *testing.T) {
	var beast ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(MutationOfCunning),
			InPlay: ct.Cards(
				ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Untamed))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Creature(
					ct.OfHouse(card.House.Untamed),
				), // a second creature, so the choice is real
			),
		},
	})

	h.P1.Play(MutationOfCunning)
	h.P1.ClickCard(beast)

	if !h.Game().HasKeyword(beast.ID(), card.Keyword.Elusive) {
		t.Error("the chosen creature should have gained elusive")
	}
	if !h.Game().HasTrait(beast.ID(), card.Traits.Mutant) {
		t.Error("the chosen creature should have gained the Mutant trait")
	}
}
