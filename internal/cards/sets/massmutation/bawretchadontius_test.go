package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bawretchadontius
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Special
//	Power:  14
//	Traits: Beast
//
//	Each friendly creature with Æmber on it gains, "Reap: Deal 4 damage to a creature."
//	Play/Fight/Reap: Exalt a friendly creature and 2 enemy creatures.
func TestBawretchadontius(t *testing.T) {
	var baw, ally, enemy1, enemy2 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(ct.Bind(&baw, Bawretchadontius), card.GiganticArt(Bawretchadontius)),
			InPlay: ct.Cards(
				ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
			),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&enemy1, ct.Creature(ct.Power(6))),
			ct.Bind(&enemy2, ct.Creature(ct.Power(6))),
		)},
	})

	// Play exalts a friendly creature and two distinct enemy creatures.
	h.P1.Play(baw)
	h.P1.ClickCard(ally)
	h.P1.ClickCard(enemy1)

	h.Expect(ally).AmberOn(1)
	h.Expect(enemy1).AmberOn(1)
	h.Expect(enemy2).AmberOn(1)

	// The ally now has Æmber, so it gains "Reap: Deal 4 damage to a creature".
	h.P1.Reap(ally)
	h.P1.ClickCard(enemy2)
	h.Expect(enemy2).Damage(4)
}
