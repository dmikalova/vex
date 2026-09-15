package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ozmo, Martianologist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Scientist
//
//	Elusive.
//	Fight/Reap: Choose one:
//	- Heal 3 damage from a Mars creature
//	- Stun a Mars creature.
func TestOzmo(t *testing.T) {
	t.Run("can stun a Mars creature when it reaps", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(Ozmo)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Reap(Ozmo)
		h.P1.ClickOption("Stun a Mars creature")

		h.Expect(foe).Stunned(true)
	})
}
