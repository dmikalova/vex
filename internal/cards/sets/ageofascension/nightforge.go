package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Nightforge
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If you have not forged a key this turn, forge a key at +4 Æmber current cost -> purge Nightforge.
var Nightforge = set.New(
	"Nightforge",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "291"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.Not{Cond: card.ForgedKey{
				Player: card.Controller,
			}},
			Then: card.ForgeKey{Extra: 4},
		}),
)
