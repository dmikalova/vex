package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Titanic Bumblebird
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Beast • Insect
//
//	Play/Reap: Destroy an enemy creature -> give a friendly creature +1 power counters equal to power of creatures destroyed this way.
func TestTitanicBumblebird(t *testing.T) {
	var bird, ally, enemy ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			InPlay: ct.Cards(
				ct.Bind(&bird, TitanicBumblebird),
				ct.Bind(&ally, ct.Creature(ct.Power(4))),
			),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&enemy, ct.Creature(ct.Power(6))),
		)},
	})

	h.P1.Reap(bird)
	// The sole enemy is destroyed; its 6 power becomes counters on a friendly.
	h.P1.ClickCard(ally)

	h.Expect(enemy).At(ct.Discard)
	h.Expect(ally).Power(10)
}
