package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Quintrino Flux
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly Creature and an enemy Creature - destroy each Creature with the same power as either of the chosen Creatures.
func TestQuintrinoFlux(t *testing.T) {
	var fChosen, fShare, fSurvive, eChosen, eShare, eSurvive ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			Hand:  ct.Cards(QuintrinoFlux),
			InPlay: ct.Cards(
				ct.Bind(&fChosen, ct.Creature(ct.Power(3))),
				ct.Bind(&fShare, ct.Creature(ct.Power(3))),
				ct.Bind(&fSurvive, ct.Creature(ct.Power(6))),
			),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&eChosen, ct.Creature(ct.Power(4))),
			ct.Bind(&eShare, ct.Creature(ct.Power(4))),
			ct.Bind(&eSurvive, ct.Creature(ct.Power(7))),
		)},
	})

	h.P1.Play(QuintrinoFlux)
	h.P1.ClickCard(fChosen) // choose the friendly creature (power 3)
	h.P1.ClickCard(eChosen) // choose the enemy creature (power 4)

	for _, c := range []ct.Card{fChosen, fShare, eChosen, eShare} {
		if h.Game().InPlay(c.ID()) {
			t.Errorf("creature %d should have been destroyed", c.ID())
		}
	}
	if !h.Game().InPlay(fSurvive.ID()) {
		t.Error("the power-6 friendly matching neither chosen power should survive")
	}
	if !h.Game().InPlay(eSurvive.ID()) {
		t.Error("the power-7 enemy matching neither chosen power should survive")
	}
}
