package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rockatiel
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast • Mutant
//
//	Elusive, Hazardous 1.
//	Play/Reap: Shuffle up to 2 creatures into their owners' decks.
func TestRockatiel(t *testing.T) {
	t.Run("shuffles up to 2 chosen creatures into their owners' decks", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(Rockatiel),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Play(Rockatiel)
		h.P1.ClickCard(ally)
		h.P1.ClickCard(foe)

		h.Expect(ally).At(ct.Deck)
		h.Expect(foe).At(ct.Deck)
	})

	t.Run("may choose no creatures", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(Rockatiel)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.Power(3))),
				ct.Creature(ct.Power(3)),
			)},
		})

		h.P1.Play(Rockatiel)
		h.P1.ClickDone()

		h.Expect(foe).At(ct.PlayArea)
	})
}
