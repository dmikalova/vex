package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Hidden Stash
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Reveal your opponent's hand. Archive a card from your opponent's hand.
func TestHiddenStash(t *testing.T) {
	t.Run(
		"archives a chosen card from the opponent's hand into the caster's archives",
		func(t *testing.T) {
			var stolen, kept ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Shadows,
					Hand:  ct.Cards(HiddenStash),
				},
				P2: ct.Side{Hand: ct.Cards(
					ct.Bind(&stolen, ct.Creature()),
					ct.Bind(&kept, ct.Creature()),
				)},
			})

			h.P1.Play(HiddenStash)
			h.P1.ClickCard(stolen)

			h.Expect(stolen).At(ct.Archives)
			h.Expect(kept).At(ct.Hand)
			if !containsCard(h.Game().Archives(0), stolen) {
				t.Error("the archived card should sit in the caster's own archives")
			}
		},
	)
}

// containsCard reports whether ids contains the card's engine id.
func containsCard(ids []engine.LocalID, c ct.Card) bool {
	for _, id := range ids {
		if id == c.ID() {
			return true
		}
	}
	return false
}
