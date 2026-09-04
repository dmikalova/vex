package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Epic Quest
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Quest
//
//	Versatile.
//	Play: Archive each friendly Knight trait creature from play.
//	Action: If you have played 7 or more Sanctum cards this turn, destroy Epic Quest, and forge a key at no cost.
func TestEpicQuest(t *testing.T) {
	t.Run("archives each friendly Knight creature in play", func(t *testing.T) {
		var knight, cleric, enemyKnight ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(EpicQuest),
				InPlay: ct.Cards(
					ct.Bind(
						&knight,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Knight)),
					),
					ct.Bind(
						&cleric,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Cleric)),
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(
					&enemyKnight,
					ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Knight)),
				),
			)},
		})

		h.P1.Play(EpicQuest)

		h.Expect(knight).At(ct.Archives)
		h.Expect(cleric).At(ct.PlayArea)
		h.Expect(enemyKnight).At(ct.PlayArea)
	})

	t.Run(
		"forges a key for free after seven Sanctum cards have been played this turn",
		func(t *testing.T) {
			var quest ct.Card
			a1 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			a2 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			a3 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			a4 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			a5 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			a6 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			a7 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Sanctum,
					InPlay: ct.Cards(ct.Bind(&quest, EpicQuest)),
					Hand:   ct.Cards(a1, a2, a3, a4, a5, a6, a7),
				},
			})

			h.P1.Play(a1)
			h.P1.Play(a2)
			h.P1.Play(a3)
			h.P1.Play(a4)
			h.P1.Play(a5)
			h.P1.Play(a6)
			h.P1.Play(a7)

			h.P1.UseAction(quest)

			h.P1.ExpectKeys(1)
			h.P1.ExpectAmber(0)
			h.Expect(quest).At(ct.Discard)
		},
	)

	t.Run("does not forge or destroy itself before seven Sanctum cards", func(t *testing.T) {
		var quest ct.Card
		a1 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
		a2 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
		a3 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
		a4 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
		a5 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
		a6 := ct.Tactic(ct.OfHouse(card.House.Sanctum))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&quest, EpicQuest)),
				Hand:   ct.Cards(a1, a2, a3, a4, a5, a6),
			},
		})

		h.P1.Play(a1)
		h.P1.Play(a2)
		h.P1.Play(a3)
		h.P1.Play(a4)
		h.P1.Play(a5)
		h.P1.Play(a6)

		h.P1.UseAction(quest)

		h.P1.ExpectKeys(0)
		h.Expect(quest).At(ct.PlayArea)
	})
}
