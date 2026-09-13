package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wretched Doll
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Play: Put a doom counter on a Creature.
//	Action: Destroy each Creature with a doom counter. Put a doom counter on a Creature.
func TestWretchedDoll(t *testing.T) {
	t.Run("play puts a doom counter on a creature", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(WretchedDoll)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&target, ct.Creature(ct.Power(4))),
			)},
		})

		h.P1.Play(WretchedDoll)

		h.Expect(target).At(ct.PlayArea)
		if h.Game().CountersOn(target.ID(), card.Counter.Doom) == 0 {
			t.Error("the chosen creature should carry a doom counter")
		}
	})

	t.Run("action destroys doomed creatures and marks another", func(t *testing.T) {
		var doomed, survivor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(WretchedDoll)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&doomed, ct.Creature(ct.Power(4))),
				ct.Bind(&survivor, ct.Creature(ct.Power(4))),
			)},
		})
		h.Game().PlaceCounter(doomed.ID(), card.Counter.Doom, 1)

		// The sweep destroys the doomed creature, then the final doom counter lands
		// on the sole survivor with no further choice.
		h.P1.UseAction(WretchedDoll)

		h.Expect(doomed).At(ct.Discard)
		h.Expect(survivor).At(ct.PlayArea)
		if h.Game().CountersOn(survivor.ID(), card.Counter.Doom) == 0 {
			t.Error("the surviving creature should carry the fresh doom counter")
		}
	})
}
