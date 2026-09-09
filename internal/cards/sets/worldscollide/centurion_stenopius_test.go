package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Centurion Stenopius
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Armor:  2
//	Traits: Dinosaur • Soldier
//
//	Centurion Stenopius gains +3 power for each Æmber on it.
//	Play/Fight/Reap: You may exalt Centurion Stenopius.
func TestCenturionStenopius(t *testing.T) {
	t.Run("exalting on play grows its power by 3 per Æmber", func(t *testing.T) {
		var centurion ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&centurion, CenturionStenopius)),
			},
		})

		h.P1.Play(CenturionStenopius)
		h.P1.ClickCard(centurion)

		h.Expect(centurion).AmberOn(1)
		h.Expect(centurion).Power(6)
	})

	t.Run("power stays at base with no Æmber on it", func(t *testing.T) {
		var centurion ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&centurion, CenturionStenopius)),
			},
		})

		h.Expect(centurion).Power(3)
	})
}
