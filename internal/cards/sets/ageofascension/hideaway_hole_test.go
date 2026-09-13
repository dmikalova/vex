package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hideaway Hole
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Versatile.
//	Action: Destroy Hideaway Hole. Each friendly Creature gains elusive until the start of your next turn.
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
}
