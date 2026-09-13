package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Techivore Pulpate
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Jelly
//
//	After a player chooses an active house, destroy each Artifact of that house.
func TestTechivorePulpate(t *testing.T) {
	t.Run("destroys each artifact of the chosen active house", func(t *testing.T) {
		var doomed, spared ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(TechivorePulpate),
			},
			P2: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&doomed, ct.Artifact(ct.OfHouse(card.House.Logos))),
					ct.Bind(&spared, ct.Artifact(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)

		h.Expect(doomed).At(ct.Discard)
		h.Expect(spared).At(ct.PlayArea)
	})

	t.Run("leaves creatures of that house untouched", func(t *testing.T) {
		var creature ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(TechivorePulpate),
			},
			P2: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&creature, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)

		h.Expect(creature).At(ct.PlayArea)
	})
}
