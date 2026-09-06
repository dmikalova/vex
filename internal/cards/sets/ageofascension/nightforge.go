package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Nightforge
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you have not forged a key this turn, you may forge a key at +4 Æmber current cost.
var Nightforge = card.New(
	"Nightforge",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 291),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.ForgedKey{
				Player: card.Controller,
				Not:    true,
			},
			Then: card.May{
				Do: card.ForgeKey{Extra: 4},
			},
		}),
)
