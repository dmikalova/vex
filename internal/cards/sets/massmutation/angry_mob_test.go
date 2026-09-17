package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Angry Mob
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Before Fight: You may discard cards from the top of your deck until you discard an Angry Mob or run out of cards -> put the discarded creature into your hand.
func TestAngryMob(t *testing.T) {
	t.Run("digs to another Angry Mob and takes it into hand", func(t *testing.T) {
		var mob, foe, skipped, found, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&mob, AngryMob)),
				Deck: ct.Cards(
					ct.Bind(&skipped, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&found, AngryMob),
					ct.Bind(&buried, ct.Creature()),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(mob, foe)
		h.P1.ClickOption("Yes")

		h.Expect(found).At(ct.Hand)
		h.Expect(skipped).At(ct.Discard)
		h.Expect(buried).At(ct.Deck)
	})

	t.Run("declining the dig leaves the deck intact", func(t *testing.T) {
		var mob, foe, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&mob, AngryMob)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(mob, foe)
		h.P1.ClickOption("No")

		h.Expect(top).At(ct.Deck)
	})
}
