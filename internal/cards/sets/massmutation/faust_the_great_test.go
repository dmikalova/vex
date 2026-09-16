package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Faust the Great
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Dinosaur
//
//	Your opponent's keys cost +1 Æmber for each friendly creature with Æmber on it.
//	Play: You may exalt a friendly creature.
func TestFaustTheGreat(t *testing.T) {
	t.Run("play exalts a friendly creature, raising the opponent's key cost", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
				Hand: ct.Cards(FaustTheGreat),
			},
		})

		h.P1.Play(FaustTheGreat)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(ally)

		h.Expect(ally).AmberOn(1)
		if got := h.Game().CurrentKeyCost(1); got != 7 {
			t.Errorf("opponent key cost = %d, want 7 (one friendly creature with Æmber)", got)
		}
		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("controller key cost = %d, want 6 (only the opponent's keys change)", got)
		}
	})

	t.Run(
		"no friendly creature holds Æmber leaves the opponent's key cost unchanged",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Saurian,
					InPlay: ct.Cards(FaustTheGreat),
				},
			})

			if got := h.Game().CurrentKeyCost(1); got != 6 {
				t.Errorf("opponent key cost = %d, want 6 (no friendly creature with Æmber)", got)
			}
		},
	)
}
