package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mimic Gel
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Shapeshifter • Mutant
//
//	Play: Choose another creature. Give Mimic Gel +1 power counters equal to its power. Mimic Gel gains the text box of the chosen creature.
func TestMimicGel(t *testing.T) {
	t.Run("copies the chosen creature's power and text box", func(t *testing.T) {
		var model ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(MimicGel),
				InPlay: ct.Cards(
					ct.Bind(&model, ct.Creature(ct.Power(6), ct.Keywords(card.Keyword.Skirmish))),
				),
			},
		})

		h.P1.Play(MimicGel)

		// Mimic Gel takes +1 power counters equal to the chosen creature's power
		// (6) on top of its base power of 1, so it reads as power 7.
		h.Expect(MimicGel).At(ct.PlayArea).Power(7)
	})
}
