package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Hideaway Hole
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Versatile.
//	Action: Destroy Hideaway Hole. Each friendly creature gains elusive until the start of your next turn.
func TestHideawayHole(t *testing.T) {
	var friend ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Shadows,
			InPlay: ct.Cards(
				HideawayHole,
				ct.Bind(&friend, ct.Creature(ct.Power(3))),
			),
		},
	})

	h.P1.UseAction(HideawayHole)

	h.Expect(HideawayHole).At(ct.Discard)
	if !h.Game().HasKeyword(friend.ID(), card.Keyword.Elusive) {
		t.Error("the friendly creature should have gained elusive")
	}

	// The grant lasts through the opponent's whole turn and lifts only at the start
	// of the controller's next turn.
	h.P1.EndTurn()
	if !h.Game().HasKeyword(friend.ID(), card.Keyword.Elusive) {
		t.Error("elusive should survive the opponent's turn")
	}
	h.P2.EndTurn()
	if h.Game().HasKeyword(friend.ID(), card.Keyword.Elusive) {
		t.Error("elusive should lift at the start of the controller's next turn")
	}
}
