package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Crassosaurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Dinosaur • Politician
//
//	Elusive.
//	Play: Crassosaurus captures 10 Æmber from any combination of players. If there are fewer than 10 Æmber on it, purge Crassosaurus.
func TestCrassosaurus(t *testing.T) {
	t.Run("capturing a full 10 keeps Crassosaurus in play", func(t *testing.T) {
		var crass ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&crass, Crassosaurus)),
			},
			P2: ct.Side{Amber: 10},
		})

		// Only the opponent holds Æmber, so all 10 come from their pool with
		// no split prompt, leaving 10 on Crassosaurus — so it is not purged.
		h.P1.Play(Crassosaurus)

		h.Expect(crass).AmberOn(10).At(ct.PlayArea)
		h.P2.ExpectAmber(0)
	})

	t.Run("capturing fewer than 10 purges Crassosaurus", func(t *testing.T) {
		var crass ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&crass, Crassosaurus)),
			},
			P2: ct.Side{Amber: 3},
		})

		// The pools hold only 3 Æmber between them, so Crassosaurus ends with
		// fewer than 10 on it and purges itself. The captured Æmber returns to
		// its owner's pool as Crassosaurus leaves play.
		h.P1.Play(Crassosaurus)

		h.Expect(crass).At(ct.Purge)
		h.P2.ExpectAmber(3)
	})
}
