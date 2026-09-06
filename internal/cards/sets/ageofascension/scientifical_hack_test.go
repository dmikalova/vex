package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Scientifical Hack
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Equation
//
//	Versatile.
//	Action: Destroy Scientifical Hack. For the remainder of the turn, you may use friendly artifacts as if they belonged to the active house.
func TestScientificalHack(t *testing.T) {
	offRelic := engine.NewCard(
		"Off Relic",
		engine.Mars,
		engine.Artifact,
		engine.Common,
		engine.WithAbility(
			engine.TriggerAction,
			engine.GainAember{Player: engine.Controller, Amount: 1},
		),
	)

	var hack, relic ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			InPlay: ct.Cards(
				ct.Bind(&hack, ScientificalHack),
				ct.Bind(&relic, offRelic),
			),
		},
	})

	h.P1.ExpectCannotUse(relic)

	h.P1.UseAction(hack)
	h.Expect(hack).At(ct.Discard)

	h.P1.UseAction(relic)
	h.P1.ExpectAmber(1)
}
