package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Shadow of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Until the start of your next turn, enemy creatures' text boxes are considered blank (except for traits).
func TestShadowOfDis(t *testing.T) {
	t.Run("blanks enemy creatures, stripping their keywords", func(t *testing.T) {
		var taunter ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(ShadowOfDis)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&taunter, ct.Creature(
					ct.Power(4), ct.Keywords(card.Keyword.Taunt))),
			)},
		})

		if !h.Game().HasKeyword(taunter.ID(), engine.Taunt) {
			t.Fatal("precondition: enemy creature should have taunt before Shadow of Dis")
		}

		h.P1.Play(ShadowOfDis)

		if h.Game().HasKeyword(taunter.ID(), engine.Taunt) {
			t.Error("a blanked enemy creature should lose its taunt keyword")
		}
	})
}
