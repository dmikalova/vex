package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Prescriptive Grammarbot
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Robot
//
//	Taunt, Hazardous 3.
//	Reap: Enrage a creature.
func TestPrescriptiveGrammarbot(t *testing.T) {
	t.Run("enrages a chosen creature when it reaps", func(t *testing.T) {
		var bot, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&bot, PrescriptiveGrammarbot)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})
		bot.Ready()

		h.P1.Reap(bot)
		h.P1.ClickCard(foe)

		if !h.Game().Enraged(foe.ID()) {
			t.Errorf("%s should be enraged", foe.Name())
		}
	})
}
