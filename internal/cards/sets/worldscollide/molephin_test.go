package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Molephin
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Hazardous 3.
//	After Æmber is stolen from you, for each Æmber stolen, deal 1 damage to each enemy creature.
func TestMolephin(t *testing.T) {
	var molephin, enemy ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Brobnar,
			InPlay: ct.Cards(ct.Bind(&enemy, Alaka)),
		},
		P2: ct.Side{
			House:  card.House.Untamed,
			InPlay: ct.Cards(ct.Bind(&molephin, Molephin)),
			Amber:  5,
		},
	})

	// Player 1 (active) steals 2 Æmber from player 2, who controls Molephin, so
	// each of Molephin's enemy creatures — player 1's — takes 2 damage.
	engine.StealAember{Amount: 2}.Resolve(
		&engine.EffectContext{
			Resolver:   h.Game(),
			Controller: 0,
		},
	)

	h.Expect(enemy).Damage(2)
	h.P2.ExpectAmber(3)
}
