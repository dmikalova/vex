package ageofascension_test

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/cards/sets/ageofascension"
)

// Rustgnawer
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast • Insect
//
//	Fight: Destroy an artifact. Resolve that card's bonus icons.
func TestRustgnawer(t *testing.T) {
	t.Run("destroys an artifact and resolves its Æmber bonus icons", func(t *testing.T) {
		var gnawer, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Amber:  0,
				InPlay: ct.Cards(ct.Bind(&gnawer, ageofascension.Rustgnawer)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(3))),
					ct.Artifact(ct.AemberBonus(2)),
				),
			},
		})

		h.P1.Fight(gnawer, foe)

		h.P1.ExpectAmber(2)
	})

	t.Run("gains nothing when the artifact has no bonus icons", func(t *testing.T) {
		var gnawer, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Amber:  0,
				InPlay: ct.Cards(ct.Bind(&gnawer, ageofascension.Rustgnawer)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(3))),
					ct.Artifact(),
				),
			},
		})

		h.P1.Fight(gnawer, foe)

		h.P1.ExpectAmber(0)
	})

	t.Run("resolves non-Æmber icons too — a Draw icon draws", func(t *testing.T) {
		var gnawer, foe, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&gnawer, ageofascension.Rustgnawer)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(3))),
					ct.Artifact(ct.Bonus(card.Bonus.Draw)),
				),
			},
		})

		h.P1.Fight(gnawer, foe)

		h.Expect(top).At(ct.Hand)
	})
}
