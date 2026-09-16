package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Maleficorn
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Mutant
//
//	After you resolve a Damage bonus icon, deal 1 damage to the same creature.
//	Enhance Damage Damage Damage Damage.
func TestMaleficorn(t *testing.T) {
	t.Run("adds 1 damage to the creature a Damage bonus icon hit", func(t *testing.T) {
		var bearer, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Maleficorn),
				Hand: ct.Cards(ct.Bind(&bearer,
					ct.Creature(ct.OfHouse(card.House.Dis), ct.Bonus(card.Bonus.Damage)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5))))},
		})

		h.P1.Play(bearer)
		h.P1.ClickCard(foe) // the Damage bonus icon hits foe

		h.Expect(foe).Damage(2) // 1 from the icon + 1 from Maleficorn
	})

	t.Run("fizzles when the bonus damage already destroyed the creature", func(t *testing.T) {
		var bearer, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Maleficorn),
				Hand: ct.Cards(ct.Bind(&bearer,
					ct.Creature(ct.OfHouse(card.House.Dis), ct.Bonus(card.Bonus.Damage)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})

		h.P1.Play(bearer)
		h.P1.ClickCard(foe) // 1 damage destroys the 1-power creature; nothing left to hit

		h.Expect(foe).At(ct.Discard)
	})
}
