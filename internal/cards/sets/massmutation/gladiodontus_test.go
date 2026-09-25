package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Gladiodontus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  15
//	Traits: Mutant
//
//	Gladiodontus deals 5 damage when fighting.
//	Gladiodontus enters play stunned.
//	Fight/Reap: If this is the first time Gladiodontus has been used this turn, ready and enrage Gladiodontus.
func TestGladiodontus(t *testing.T) {
	t.Run("its first reap readies and enrages it", func(t *testing.T) {
		var glad ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&glad, Gladiodontus)),
			},
		})

		h.P1.Reap(glad)

		h.Expect(glad).Ready().Enraged(true)
	})

	t.Run("its second use does not ready it again", func(t *testing.T) {
		var glad ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&glad, Gladiodontus)),
			},
		})

		h.P1.Reap(glad) // first use: readies it
		h.P1.Reap(glad) // second use: condition not met

		h.Expect(glad).Exhausted()
	})
}
