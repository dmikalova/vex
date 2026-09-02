package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mushroom Man
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Fungus • Human
//
//	Mushroom Man gains +3 power for each unforged key you have.
func TestMushroomMan(t *testing.T) {
	t.Run("gets +3 power per unforged key", func(t *testing.T) {
		var man ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&man, MushroomMan)),
			},
		})

		h.Expect(man).Power(11) // 2 + 3 per unforged key, none forged
	})

	t.Run("shrinks as its controller forges keys", func(t *testing.T) {
		var man ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Keys:   2,
				InPlay: ct.Cards(ct.Bind(&man, MushroomMan)),
			},
		})

		h.Expect(man).Power(5)
	})
}
