package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mutation of Instinct
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Choose a creature. The chosen creature gains skirmish and the Mutant trait until the start of your next turn.
func TestMutationOfInstinct(t *testing.T) {
	var beast ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(MutationOfInstinct),
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

	h.P1.Play(MutationOfInstinct)
	h.P1.ClickCard(beast)

	if !h.Game().HasKeyword(beast.ID(), card.Keyword.Skirmish) {
		t.Error("the chosen creature should have gained skirmish")
	}
	if !h.Game().HasTrait(beast.ID(), card.Traits.Mutant) {
		t.Error("the chosen creature should have gained the Mutant trait")
	}
}
