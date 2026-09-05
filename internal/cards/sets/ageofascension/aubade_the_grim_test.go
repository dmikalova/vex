package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Aubade the Grim
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Spirit • Knight
//
//	Play: Aubade the Grim captures 3 Æmber from your opponent.
//	Reap: Move 1 Æmber from Aubade the Grim to the common supply.
func TestAubadeTheGrim(t *testing.T) {
	t.Run("captures 3 Æmber when played", func(t *testing.T) {
		var aubade ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(ct.Bind(&aubade, AubadeTheGrim))},
			P2: ct.Side{Amber: 4},
		})

		h.P1.Play(AubadeTheGrim)

		h.Expect(aubade).AmberOn(3)
	})

	t.Run("discards 1 of its captured Æmber when it reaps", func(t *testing.T) {
		var aubade ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&aubade, AubadeTheGrim)),
			},
			P2: ct.Side{},
		})
		h.Game().AddAmberOn(aubade.ID(), 3)

		h.P1.Reap(aubade)

		h.Expect(aubade).AmberOn(2)
	})
}
