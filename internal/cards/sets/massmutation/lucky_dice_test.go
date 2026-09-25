package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Lucky Dice
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Lucky Dice. During your opponent's next turn, each friendly creature cannot be dealt damage.
func TestLuckyDice(t *testing.T) {
	t.Run("protects friendly creatures only during the opponent's next turn", func(t *testing.T) {
		var mine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					LuckyDice,
					ct.Bind(&mine, ct.Creature(ct.Power(5))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar)))},
		})
		g := h.Game()

		h.P1.UseAction(LuckyDice)
		h.Expect(LuckyDice).At(ct.Discard)
		// Dormant on the caster's own turn.
		if g.DamageImmune(mine.ID()) {
			t.Error("immunity should be dormant during the caster's turn")
		}

		// During the opponent's turn, the friendly creature cannot be dealt damage.
		g.EndPlayPhase(0)
		g.StartTurn(1)
		if err := g.ChooseHouse(1, card.House.Brobnar); err != nil {
			t.Fatalf("ChooseHouse: %v", err)
		}
		g.DealDamage(1, []engine.DamageTarget{{ID: mine.ID(), Amount: 3}})
		if g.Damage(mine.ID()) != 0 {
			t.Errorf(
				"friendly creature took %d damage during opponent's turn, want 0",
				g.Damage(mine.ID()),
			)
		}

		// The protection lifts once the opponent's turn ends.
		g.EndPlayPhase(1)
		g.StartTurn(0)
		if err := g.ChooseHouse(0, card.House.Shadows); err != nil {
			t.Fatalf("ChooseHouse: %v", err)
		}
		g.DealDamage(0, []engine.DamageTarget{{ID: mine.ID(), Amount: 3}})
		if g.Damage(mine.ID()) != 3 {
			t.Errorf(
				"friendly creature took %d after protection lifted, want 3",
				g.Damage(mine.ID()),
			)
		}
	})
}
