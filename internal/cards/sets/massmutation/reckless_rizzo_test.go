package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Reckless Rizzo
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 2 Æmber. Until the start of your next turn, Reckless Rizzo loses elusive.
func TestRecklessRizzo(t *testing.T) {
	var rizzo ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Shadows,
			InPlay: ct.Cards(ct.Bind(&rizzo, RecklessRizzo)),
		},
		P2: ct.Side{Amber: 2},
	})

	if !h.Game().HasKeyword(rizzo.ID(), engine.Elusive) {
		t.Fatal("Reckless Rizzo should start elusive")
	}

	h.P1.UseAction(rizzo)

	h.P1.ExpectAmber(2)
	h.P2.ExpectAmber(0)
	if h.Game().HasKeyword(rizzo.ID(), engine.Elusive) {
		t.Error("Reckless Rizzo should have lost elusive after acting")
	}

	// The loss lasts through the opponent's whole turn and lifts only at the start
	// of Reckless Rizzo's controller's next turn.
	h.P1.EndTurn()
	if h.Game().HasKeyword(rizzo.ID(), engine.Elusive) {
		t.Error("the loss should survive the opponent's turn")
	}
	h.P2.EndTurn()
	if !h.Game().HasKeyword(rizzo.ID(), engine.Elusive) {
		t.Error("elusive should return at the start of Reckless Rizzo's controller's next turn")
	}
}
