package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Curse of Vanity
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Exalt a friendly creature and an enemy creature.
func TestCurseOfVanity(t *testing.T) {
	t.Run("exalts a friendly and an enemy creature", func(t *testing.T) {
		var friendly, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(CurseOfVanity),
				InPlay: ct.Cards(ct.Bind(&friendly, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature()))},
		})

		h.P1.Play(CurseOfVanity)

		h.Expect(friendly).AmberOn(1)
		h.Expect(enemy).AmberOn(1)
	})
}
