package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Knowledge is Power
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose one:
//	- Archive a card from your hand
//	- For each card in your archives, gain 1 Æmber.
func TestKnowledgeIsPower(t *testing.T) {
	t.Run("gains 1 Æmber for each archived card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(KnowledgeIsPower),
				Archives: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(1)),
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(1)),
				),
			},
		})

		h.P1.Play(KnowledgeIsPower)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("gain 1") // the "for each archived card, gain" option

		h.P1.ExpectAmber(2) // 1 per each of the 2 archived cards
	})
}
