package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Champion's Challenge
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy each enemy creature except the most powerful enemy creature and each friendly creature except the most powerful friendly creature, and ready and fight with a friendly creature.
func TestChampionsChallenge(t *testing.T) {
	t.Run("wipes each side to its strongest, then the ally fights", func(t *testing.T) {
		var strongAlly, weakAlly, strongFoe, weakFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ChampionsChallenge),
				InPlay: ct.Cards(
					ct.Bind(&strongAlly, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(8))),
					ct.Bind(&weakAlly, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&strongFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(7))),
					ct.Bind(&weakFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
		})

		h.P1.Play(ChampionsChallenge)

		h.Expect(weakAlly).At(ct.Discard)
		h.Expect(weakFoe).At(ct.Discard)
		h.Expect(strongFoe).At(ct.Discard)   // the surviving ally (8) fights and destroys it (7)
		h.Expect(strongAlly).At(ct.PlayArea) // and survives the exchange
	})
}
