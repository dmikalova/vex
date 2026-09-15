package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Quixxle Stone
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	If a player has more creatures in play than their opponent, they cannot play creatures.
func TestQuixxleStone(t *testing.T) {
	t.Run("bars the player who controls more creatures from playing creatures", func(t *testing.T) {
		var mine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					QuixxleStone,
					ct.Creature(ct.OfHouse(card.House.StarAlliance)),
				),
				Hand: ct.Cards(ct.Bind(&mine, ct.Creature(ct.OfHouse(card.House.StarAlliance)))),
			},
		})

		h.P1.ExpectCannotPlay(mine)
	})

	t.Run("only bars creatures, not other card types", func(t *testing.T) {
		var tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					QuixxleStone,
					ct.Creature(ct.OfHouse(card.House.StarAlliance)),
				),
				Hand: ct.Cards(ct.Bind(&tactic, ct.Tactic(ct.OfHouse(card.House.StarAlliance)))),
			},
		})

		h.P1.Play(tactic)
	})

	t.Run("allows creature plays once the counts are equal", func(t *testing.T) {
		var mine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					QuixxleStone,
					ct.Creature(ct.OfHouse(card.House.StarAlliance)),
				),
				Hand: ct.Cards(ct.Bind(&mine, ct.Creature(ct.OfHouse(card.House.StarAlliance)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.StarAlliance))),
			},
		})

		h.P1.Play(mine)
	})
}
