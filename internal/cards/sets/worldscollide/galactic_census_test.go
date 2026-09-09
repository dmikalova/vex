package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Galactic Census
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are 3 or more houses represented among creatures in play, gain 1 Æmber. If there are 5 or more houses represented among creatures in play, gain 1 Æmber. If there are 6 or more houses represented among creatures in play, gain 1 Æmber.
func TestGalacticCensus(t *testing.T) {
	// The played Tactic contributes its bonus Æmber before its effect resolves,
	// so every expected total is that 1 bonus plus the tier payout.
	t.Run("fewer than 3 houses pays only the bonus Æmber", func(t *testing.T) {
		var census ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				Hand:   ct.Cards(ct.Bind(&census, GalacticCensus)),
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Dis)))},
		})

		h.P1.Play(census)

		h.P1.ExpectAmber(1)
	})

	t.Run("three or four houses gains 1 Æmber", func(t *testing.T) {
		var census ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ct.Bind(&census, GalacticCensus)),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar)),
					ct.Creature(ct.OfHouse(card.House.Dis)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Logos)))},
		})

		h.P1.Play(census)

		h.P1.ExpectAmber(2)
	})

	t.Run("exactly five houses gains 2 Æmber", func(t *testing.T) {
		var census ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ct.Bind(&census, GalacticCensus)),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar)),
					ct.Creature(ct.OfHouse(card.House.Dis)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Mars)),
				ct.Creature(ct.OfHouse(card.House.Sanctum)),
			)},
		})

		h.P1.Play(census)

		h.P1.ExpectAmber(3)
	})

	t.Run("six or more houses gains 3 Æmber", func(t *testing.T) {
		var census ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(ct.Bind(&census, GalacticCensus)),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar)),
					ct.Creature(ct.OfHouse(card.House.Dis)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Mars)),
				ct.Creature(ct.OfHouse(card.House.Sanctum)),
				ct.Creature(ct.OfHouse(card.House.Shadows)),
			)},
		})

		h.P1.Play(census)

		h.P1.ExpectAmber(4)
	})
}
