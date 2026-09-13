package ageofascension_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
)

// Rustgnawer
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast • Insect
//
//	Fight: Destroy an Artifact. For each Æmber bonus on it, gain 1 Æmber.
func TestRustgnawer(t *testing.T) {
	t.Run("destroys an artifact and gains its Æmber bonus", func(t *testing.T) {
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

	t.Run("gains nothing when the artifact has no Æmber bonus", func(t *testing.T) {
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
}
