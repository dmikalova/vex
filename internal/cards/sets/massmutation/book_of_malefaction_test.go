package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Book of Malefaction
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item • Law
//
//	Versatile.
//	After Æmber is stolen from you, for each Æmber stolen, put a warrant counter on Book of Malefaction.
//	Action: Remove a warrant counter from Book of Malefaction -> purge a creature.
func TestBookOfMalefaction(t *testing.T) {
	t.Run("gains a warrant counter per Æmber stolen, then spends one to purge", func(t *testing.T) {
		var book, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&book, BookOfMalefaction)),
				Amber:  3,
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		// Player 2 steals 2 Æmber from Book's controller, arming two warrant counters.
		engine.StealAember{Amount: 2}.Resolve(
			&engine.EffectContext{Resolver: h.Game(), Controller: 1},
		)

		// The Omni removes one warrant counter and purges a creature. The lone
		// enemy creature is the only legal target, so it is purged without a prompt.
		h.P1.UseAction(BookOfMalefaction)

		h.Expect(victim).At(ct.Purge)
	})

	t.Run("with no warrant counter the Omni purges nothing", func(t *testing.T) {
		var book, bystander ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&book, BookOfMalefaction)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&bystander, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(BookOfMalefaction)

		h.Expect(bystander).At(ct.PlayArea)
	})
}
