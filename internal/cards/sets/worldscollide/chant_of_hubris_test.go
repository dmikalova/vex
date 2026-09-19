package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chant of Hubris
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Move 1 Æmber from a creature to another creature.
func TestChantOfHubris(t *testing.T) {
	t.Run("moves 1 Æmber from one creature onto another", func(t *testing.T) {
		var from, onto, bystander ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ChantOfHubris),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&from, ct.Creature(ct.Power(3))),
				ct.Bind(&onto, ct.Creature(ct.Power(3))),
				ct.Bind(&bystander, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[from.ID()].Amber = 1

		h.P1.Play(ChantOfHubris)
		h.P1.ClickCard(onto)

		h.Expect(from).AmberOn(0)
		h.Expect(onto).AmberOn(1)
		h.Expect(bystander).AmberOn(0)
	})

	// "another creature" excludes the creature the Æmber is leaving, so with only
	// one other creature in play there is nothing left to choose between: moving
	// Æmber off a card and back onto it is not a move.
	t.Run("never offers the creature the Æmber is leaving", func(t *testing.T) {
		var from, onto ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ChantOfHubris),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&from, ct.Creature(ct.Power(3))),
				ct.Bind(&onto, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[from.ID()].Amber = 1

		h.P1.Play(ChantOfHubris)

		h.Expect(from).AmberOn(0)
		h.Expect(onto).AmberOn(1)
	})
}
