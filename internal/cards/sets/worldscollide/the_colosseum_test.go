package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Colosseum
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Versatile.
//	After an enemy creature is destroyed while fighting, put a glory counter on The Colosseum.
//	Action: If there are 6 or more glory counters on The Colosseum, remove 6 glory counters from The Colosseum, and forge a key at current cost -> purge The Colosseum.
func TestTheColosseum(t *testing.T) {
	t.Run("an enemy creature destroyed while fighting adds a glory counter", func(t *testing.T) {
		var colosseum, fighter, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&colosseum, TheColosseum),
					ct.Bind(&fighter, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(10))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(fighter, enemy)

		if got := h.Game().CountersOn(colosseum.ID(), card.Counter.Glory); got != 1 {
			t.Errorf("glory counters after the kill = %d, want 1", got)
		}
	})

	t.Run("the Action forges a key once six glory counters are gathered", func(t *testing.T) {
		var colosseum ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Amber:  6,
				InPlay: ct.Cards(ct.Bind(&colosseum, TheColosseum)),
			},
		})
		h.Game().PlaceCounter(colosseum.ID(), card.Counter.Glory, 6)

		h.P1.UseAction(colosseum)

		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(0) // paid the current key cost of six
		if got := h.Game().CountersOn(colosseum.ID(), card.Counter.Glory); got != 0 {
			t.Errorf("glory counters after forging = %d, want 0", got)
		}
	})
}
