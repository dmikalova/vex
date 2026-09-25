package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Seeker of Truth
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Armor:  1
//	Traits: Human
//
//	Fight: You may use a friendly non-Sanctum creature.
func TestSeekerOfTruth(t *testing.T) {
	t.Run("with no non-Sanctum friendly creature the fight just resolves", func(t *testing.T) {
		var seeker, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&seeker, SeekerOfTruth)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Fight(seeker, foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(seeker).Damage(1) // took 2 from the fight, armor 1 absorbs 1
	})

	t.Run("uses a friendly non-Sanctum creature", func(t *testing.T) {
		var seeker, foe, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(
				ct.Bind(&seeker, SeekerOfTruth),
				ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Logos))),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Fight(seeker, foe)
		h.P1.ClickCard(other) // choose the non-Sanctum creature to use (auto-reaps)

		h.Expect(other).Exhausted()
	})
}
