package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sci. Officer Morpheus
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Shapeshifter • Scientist
//
//	After a creature is played, if it is a friendly creature, trigger the play effect of it.
func TestSciOfficerMorpheus(t *testing.T) {
	t.Run("re-triggers a friendly creature's play effect", func(t *testing.T) {
		var dux ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Amber:  6, // enough pool to exalt four times
				InPlay: ct.Cards(SciOfficerMorpheus),
				Hand:   ct.Cards(ct.Bind(&dux, GrimlocusDux)),
			},
		})

		h.P1.Play(dux) // natural Play exalts 2; Morpheus re-triggers for 2 more

		h.Expect(dux).AmberOn(4)
	})

	t.Run("a friendly creature with no play effect is a no-op", func(t *testing.T) {
		var vanilla ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(SciOfficerMorpheus),
				Hand: ct.Cards(ct.Bind(&vanilla,
					ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3)))),
			},
		})

		h.P1.Play(vanilla)

		h.Expect(vanilla).At(ct.PlayArea).AmberOn(0)
	})
}
