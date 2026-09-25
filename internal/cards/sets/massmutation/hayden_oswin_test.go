package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Hayden Oswin
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human
//
//	Reap: For each upgrade on Hayden Oswin, gain 1 Æmber.
func TestHaydenOswin(t *testing.T) {
	t.Run("gains only the base reap Æmber with no upgrades", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(HaydenOswin),
			},
		})

		h.P1.Reap(HaydenOswin)

		h.P1.ExpectAmber(1) // 1 for reaping, 0 upgrades
	})

	t.Run("gains 1 Æmber for each upgrade on it", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(HaydenOswin, ct.Upgrade(), ct.Upgrade()),
				),
			},
		})

		h.P1.Reap(HaydenOswin)

		h.P1.ExpectAmber(3) // 1 for reaping, 2 for the upgrades
	})
}
