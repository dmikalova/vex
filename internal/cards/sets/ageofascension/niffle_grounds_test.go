package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Niffle Grounds
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Action: Choose a creature. For the remainder of the turn, it loses taunt and elusive.
func TestNiffleGrounds(t *testing.T) {
	var target ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Untamed,
			InPlay: ct.Cards(NiffleGrounds),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&target, ct.Creature(
					ct.Power(3),
					ct.Keywords(card.Keyword.Taunt, card.Keyword.Elusive),
				)),
			),
		},
	})

	if !h.Game().HasKeyword(target.ID(), card.Keyword.Taunt) ||
		!h.Game().HasKeyword(target.ID(), card.Keyword.Elusive) {
		t.Fatal("the target should start with taunt and elusive")
	}

	h.P1.UseAction(NiffleGrounds)

	if h.Game().HasKeyword(target.ID(), card.Keyword.Taunt) {
		t.Error("taunt should be lost for the remainder of the turn")
	}
	if h.Game().HasKeyword(target.ID(), card.Keyword.Elusive) {
		t.Error("elusive should be lost for the remainder of the turn")
	}
}
