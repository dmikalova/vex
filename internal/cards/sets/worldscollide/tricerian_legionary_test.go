package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tricerian Legionary
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Taunt.
//	Play: Ward a friendly Creature.
func TestTricerianLegionary(t *testing.T) {
	t.Run("wards a chosen friendly creature when played", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(TricerianLegionary),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(3)))),
			},
			P2: ct.Side{},
		})

		h.P1.Play(TricerianLegionary)
		h.P1.ClickCard(ally)

		if !h.Game().Warded(ally.ID()) {
			t.Errorf("%s should be warded", ally.Name())
		}
	})
}
