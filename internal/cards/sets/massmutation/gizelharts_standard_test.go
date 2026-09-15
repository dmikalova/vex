package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gizelhart's Standard
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Each friendly creature with Æmber on it gains +1 armor.
//	Play: Exalt a friendly creature.
func TestGizelhartsStandard(t *testing.T) {
	t.Run("grants +1 armor only to friendly creatures with Æmber on them", func(t *testing.T) {
		var withAember, bare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					GizelhartsStandard,
					ct.Bind(&withAember, ct.Creature(ct.Armor(1))),
					ct.Bind(&bare, ct.Creature(ct.Armor(1))),
				),
			},
		})
		h.Game().State.Cards[withAember.ID()].Amber = 1

		h.Expect(withAember).Armor(2)
		h.Expect(bare).Armor(1)
	})

	t.Run("play exalts a friendly creature", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Hand:   ct.Cards(GizelhartsStandard),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature())),
			},
		})

		h.P1.Play(GizelhartsStandard)

		h.Expect(ally).AmberOn(1)
	})
}
