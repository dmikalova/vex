package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Spoils of Battle
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: A friendly creature captures 1 Æmber from your opponent. Each creature with Æmber on it captures 1 Æmber from its opponent.
func TestSpoilsOfBattle(t *testing.T) {
	t.Run(
		"a friendly creature captures, then each Æmber-bearing creature captures from its opponent",
		func(t *testing.T) {
			var friendly, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Saurian,
					Hand:   ct.Cards(SpoilsOfBattle),
					InPlay: ct.Cards(ct.Bind(&friendly, ct.Creature(ct.Power(3)))),
					Amber:  5,
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3)))),
					Amber:  5,
				},
			})
			h.Game().State.Cards[enemy.ID()].Amber = 1

			h.P1.Play(SpoilsOfBattle)

			// Sentence 1: the lone friendly creature captures 1 from the opponent's pool.
			// Sentence 2: it (now Æmber-bearing) captures 1 more from the opponent, and
			// the enemy creature captures 1 from its own opponent (P1).
			h.Expect(friendly).AmberOn(2)
			h.Expect(enemy).AmberOn(2)
			h.P2.ExpectAmber(3) // 5 - 1 (sentence 1) - 1 (friendly in sentence 2)
		},
	)
}
