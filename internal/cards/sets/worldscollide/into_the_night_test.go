package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Into the Night
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Until the start of your next turn, non-Shadows Creatures cannot be used to fight.
func TestIntoTheNight(t *testing.T) {
	t.Run("spares Shadows creatures and bars every other house", func(t *testing.T) {
		var shadowAlly, foe, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(IntoTheNight),
				InPlay: ct.Cards(
					ct.Bind(&shadowAlly, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(IntoTheNight)
		// On the caster's own Shadows turn, a spared Shadows creature still fights.
		h.P1.Fight(shadowAlly, foe)
		h.P1.EndTurn()

		// On the opponent's turn, their non-Shadows creatures cannot fight.
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.ExpectCannotUseTo(enemy, engine.FightUse)
	})

	t.Run("lifts once the caster's next turn begins", func(t *testing.T) {
		var enemy, target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(IntoTheNight),
				InPlay: ct.Cards(
					ct.Bind(&target, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(IntoTheNight)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.ExpectCannotUseTo(enemy, engine.FightUse)
		h.P2.EndTurn()

		// The caster's next turn: the bar is gone, so its opponent's creature could
		// fight again — shown here by the opponent readying and fighting on the turn
		// after next.
		h.P1.ChooseHouse(card.House.Shadows)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Fight(enemy, target)
	})
}
